package recipe

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func globPaths(pattern string) ([]string, error) {
	pattern = os.ExpandEnv(pattern)
	if !strings.Contains(pattern, "**") {
		return filepath.Glob(pattern)
	}

	root, leafPattern, err := splitRecursiveGlob(pattern)
	if err != nil {
		return nil, err
	}

	return recursiveDirGlob(root, leafPattern)
}

func splitRecursiveGlob(pattern string) (string, string, error) {
	separatorGlobstarSeparator := string(filepath.Separator) + "**" + string(filepath.Separator)
	parts := strings.Split(pattern, separatorGlobstarSeparator)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("recursive glob %q must use exactly one %q segment", pattern, "**")
	}
	if strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return "", "", fmt.Errorf("recursive glob %q must include a root and a leaf pattern", pattern)
	}
	if strings.Contains(parts[1], string(filepath.Separator)) {
		return "", "", fmt.Errorf("recursive glob %q only supports a leaf pattern after %q", pattern, "**")
	}

	return parts[0], parts[1], nil
}

func recursiveDirGlob(root, leafPattern string) ([]string, error) {
	var matches []string

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !entry.IsDir() {
			return nil
		}

		matched, err := filepath.Match(leafPattern, entry.Name())
		if err != nil {
			return err
		}
		if matched {
			matches = append(matches, path)
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	slices.Sort(matches)
	return matches, nil
}
