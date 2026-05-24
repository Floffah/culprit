package shell

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/floffah/maculprit/internal/recipe"
)

func RemovesToBash(removes []recipe.Remove) string {
	var script strings.Builder

	script.WriteString("#!/usr/bin/env bash\n")
	script.WriteString("set -euo pipefail\n\n")

	var previousRecipeName string
	for _, remove := range removes {
		removeSize, err := getSize(remove.Path)
		if err != nil {
			removeSize = 0
		}

		if remove.RecipeName != previousRecipeName {
			script.WriteString("#")
			writeComment(&script, "Recipe", remove.RecipeName)
			script.WriteString("\n")
			previousRecipeName = remove.RecipeName
		}

		writeComment(&script, "Reason", remove.Reason)
		writeComment(&script, "Matcher", remove.Matcher)
		writeComment(&script, "Size", formatSize(removeSize))
		script.WriteString("rm -rf -- ")
		script.WriteString(shellQuote(remove.Path))
		script.WriteString("\n\n")
	}

	return script.String()
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
		err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
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
