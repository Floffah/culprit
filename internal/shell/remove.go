package shell

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/floffah/maculprit/internal/recipe"
)

func RemovesToBash(removes []recipe.Remove) string {
	var script string

	var previousRecipeName string
	for _, remove := range removes {
		removeSize, err := getSize(remove.Path)
		if err != nil {
			removeSize = 0
		}

		if remove.RecipeName != previousRecipeName {
			script += "# -- Recipe: " + remove.RecipeName + " --\n\n"
			previousRecipeName = remove.RecipeName
		}

		script += "# Reason: " + remove.Reason + "\n"
		script += "# Matcher: " + remove.Matcher + "\n"
		script += "# Size: " + formatSize(removeSize) + "\n"
		script += "rm -rf " + strconv.Quote(remove.Path) + "\n\n"
	}

	return script
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

func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

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
