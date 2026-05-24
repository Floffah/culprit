package clean

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/floffah/culprit/internal/recipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipeTypeFromArgsDefaultsToSoft(t *testing.T) {
	assert.Equal(t, "soft", recipeTypeFromArgs(nil))
}

func TestRecipeTypeFromArgsUsesFirstArg(t *testing.T) {
	assert.Equal(t, "hard", recipeTypeFromArgs([]string{"hard"}))
}

func TestFilterRecipesByType(t *testing.T) {
	recipes := []recipe.Recipe{
		{Name: "Soft One", Type: recipe.TypeSoft},
		{Name: "Hard One", Type: recipe.TypeHard},
		{Name: "Soft Two", Type: recipe.TypeSoft},
	}

	filtered := filterRecipesByType(recipes, "soft")

	require.Len(t, filtered, 2)
	assert.Equal(t, "Soft One", filtered[0].Name)
	assert.Equal(t, "Soft Two", filtered[1].Name)
}

func TestWriteScriptFileCreatesNewFile(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "cleanup.sh")

	require.NoError(t, writeScriptFile(outputPath, "script contents", false))

	contents, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Equal(t, "script contents", string(contents))
}

func TestWriteScriptFileRefusesExistingFileWithoutForce(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "cleanup.sh")
	require.NoError(t, os.WriteFile(outputPath, []byte("keep me"), 0644))

	err := writeScriptFile(outputPath, "replace me", false)

	require.Error(t, err)
	assert.ErrorContains(t, err, "already exists")

	contents, readErr := os.ReadFile(outputPath)
	require.NoError(t, readErr)
	assert.Equal(t, "keep me", string(contents))
}

func TestWriteScriptFileOverwritesExistingFileWithForce(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "cleanup.sh")
	require.NoError(t, os.WriteFile(outputPath, []byte("replace me"), 0644))

	require.NoError(t, writeScriptFile(outputPath, "new script", true))

	contents, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Equal(t, "new script", string(contents))
}
