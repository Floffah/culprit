package shell

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/floffah/culprit/internal/recipe"
)

type RenderOptions struct {
	IncludeSizes bool
}

func RemovesToBash(removes []recipe.Remove) string {
	return RemovesToBashWithOptions(removes, RenderOptions{IncludeSizes: true})
}

func RemovesToBashWithOptions(removes []recipe.Remove, options RenderOptions) string {
	var script strings.Builder

	script.WriteString("#!/usr/bin/env bash\n")
	script.WriteString("set -euo pipefail\n\n")

	var removeSizes map[string]int64
	if options.IncludeSizes {
		removeSizes = getRemoveSizes(removes)
	}

	var previousRecipeName string
	for _, remove := range removes {
		if remove.RecipeName != previousRecipeName {
			script.WriteString("#")
			writeComment(&script, "Recipe", remove.RecipeName)
			script.WriteString("\n")
			previousRecipeName = remove.RecipeName
		}

		writeComment(&script, "Reason", remove.Reason)
		writeComment(&script, "Matcher", remove.Matcher)
		if options.IncludeSizes {
			writeComment(&script, "Size", formatSize(removeSizes[removeSizeKey(remove)]))
		}
		if remove.Command != "" {
			script.WriteString(remove.Command)
		} else {
			script.WriteString("rm -rf -- ")
			script.WriteString(shellQuote(remove.Path))
		}
		script.WriteString("\n\n")
	}

	return script.String()
}

func getRemoveSizes(removes []recipe.Remove) map[string]int64 {
	type job struct {
		key    string
		remove recipe.Remove
	}

	jobsByKey := make(map[string]job)
	for _, remove := range removes {
		key := removeSizeKey(remove)
		if _, exists := jobsByKey[key]; exists {
			continue
		}
		jobsByKey[key] = job{key: key, remove: remove}
	}

	jobs := make(chan job)
	results := make(map[string]int64, len(jobsByKey))
	var resultsMu sync.Mutex

	workerCount := min(runtime.NumCPU(), len(jobsByKey))
	if workerCount < 1 {
		workerCount = 1
	}

	var wg sync.WaitGroup
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				size := getRemoveSize(job.remove)
				resultsMu.Lock()
				results[job.key] = size
				resultsMu.Unlock()
			}
		}()
	}

	for _, job := range jobsByKey {
		jobs <- job
	}
	close(jobs)
	wg.Wait()

	return results
}

func removeSizeKey(remove recipe.Remove) string {
	if remove.Command != "" {
		return "command:" + remove.Command + "\x00" + strings.Join(remove.RelatedPaths, "\x00")
	}

	return "path:" + remove.Path
}

func getRemoveSize(remove recipe.Remove) int64 {
	if remove.Command != "" {
		return getRelatedPathsSize(remove.RelatedPaths)
	}

	removeSize, err := getSize(remove.Path)
	if err != nil {
		return 0
	}

	return removeSize
}

func getRelatedPathsSize(paths []string) int64 {
	var totalSize int64

	for _, path := range paths {
		matches, err := filepath.Glob(os.ExpandEnv(path))
		if err != nil {
			continue
		}
		for _, match := range matches {
			size, err := getSize(match)
			if err != nil {
				continue
			}
			totalSize += size
		}
	}

	return totalSize
}

func writeComment(script *strings.Builder, label, value string) {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")

	lines := strings.Split(value, "\n")
	if len(lines) == 0 {
		script.WriteString("# ")
		script.WriteString(label)
		script.WriteString(":\n")
		return
	}

	script.WriteString("# ")
	script.WriteString(label)
	script.WriteString(": ")
	script.WriteString(lines[0])
	script.WriteString("\n")

	for _, line := range lines[1:] {
		script.WriteString("#   ")
		script.WriteString(line)
		script.WriteString("\n")
	}
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}

	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func getSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if fileInfo.IsDir() {
		var totalSize int64 = 0
		err := filepath.WalkDir(path, func(_ string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !entry.IsDir() {
				info, err := entry.Info()
				if err != nil {
					return err
				}
				totalSize += info.Size()
			}
			return nil
		})
		if err != nil {
			return 0, err
		}
		return totalSize, nil
	}
	return fileInfo.Size(), nil
}

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

func formatSize(size int64) string {
	switch {
	case size >= GB:
		return formatFloat(float64(size)/GB) + " GB"
	case size >= MB:
		return formatFloat(float64(size)/MB) + " MB"
	case size >= KB:
		return formatFloat(float64(size)/KB) + " KB"
	default:
		return formatFloat(float64(size)) + " B"
	}
}

func formatFloat(f float64) string {
	if f >= 100 {
		return strconv.FormatFloat(f, 'f', 0, 64)
	} else if f >= 10 {
		return strconv.FormatFloat(f, 'f', 1, 64)
	} else {
		return strconv.FormatFloat(f, 'f', 2, 64)
	}
}
