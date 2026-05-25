package recipe

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobPathsSupportsGlobstar(t *testing.T) {
	root := t.TempDir()
	createDir(t, filepath.Join(root, "node_modules"))
	createDir(t, filepath.Join(root, "app", "node_modules"))
	createDir(t, filepath.Join(root, "app", "node_modules", "dependency", "node_modules"))
	createDir(t, filepath.Join(root, "nested", "app", "node_modules"))
	createDir(t, filepath.Join(root, "nested", "app", "src"))

	matches, err := globPaths(filepath.Join(root, "**", "node_modules"))

	require.NoError(t, err)
	assert.Equal(t, []string{
		filepath.Join(root, "app", "node_modules"),
		filepath.Join(root, "nested", "app", "node_modules"),
		filepath.Join(root, "node_modules"),
	}, matches)
}

func TestGlobPathsKeepsStandardGlobBehavior(t *testing.T) {
	root := t.TempDir()
	createDir(t, filepath.Join(root, "cache-a"))
	createDir(t, filepath.Join(root, "cache-b"))
	createDir(t, filepath.Join(root, "other"))

	matches, err := globPaths(filepath.Join(root, "cache-*"))

	require.NoError(t, err)
	assert.Equal(t, []string{
		filepath.Join(root, "cache-a"),
		filepath.Join(root, "cache-b"),
	}, matches)
}

func TestGlobPathsSupportsLeafPatternAfterGlobstar(t *testing.T) {
	root := t.TempDir()
	createDir(t, filepath.Join(root, "app", "node_modules"))
	createDir(t, filepath.Join(root, "app", "node_cache"))
	createDir(t, filepath.Join(root, "app", "other"))

	matches, err := globPaths(filepath.Join(root, "**", "node_*"))

	require.NoError(t, err)
	assert.Equal(t, []string{
		filepath.Join(root, "app", "node_cache"),
		filepath.Join(root, "app", "node_modules"),
	}, matches)
}

func TestGlobPathsRejectsComplexGlobstar(t *testing.T) {
	_, err := globPaths(filepath.Join(t.TempDir(), "**", "node_modules", "*"))

	require.Error(t, err)
	assert.ErrorContains(t, err, "only supports a leaf pattern")
}

func TestDeriveRecipeRemovesSupportsGlobstar(t *testing.T) {
	root := t.TempDir()
	createDir(t, filepath.Join(root, "web", "site", "node_modules"))

	removes, err := DeriveRecipeRemoves(Recipe{
		Name: "Node Modules",
		Targets: []Target{
			{
				Reason: "test",
				Kind:   TargetAbsolute,
				Path:   filepath.Join(root, "**", "node_modules"),
			},
		},
	})

	require.NoError(t, err)
	require.Len(t, removes, 1)
	assert.Equal(t, filepath.Join(root, "web", "site", "node_modules"), removes[0].Path)
}

func createDir(t *testing.T, path string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(path, 0755))
}
