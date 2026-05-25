package clean

import (
	"os"
	"testing"

	"github.com/floffah/culprit/internal/recipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectRecipeInputsUsesDefaultWithoutTTY(t *testing.T) {
	valuesByRecipe, err := collectRecipeInputs([]recipe.Recipe{
		{
			Name: "Input Recipe",
			Inputs: []recipe.Input{
				{
					ID:      "project_path",
					Kind:    recipe.InputPath,
					Prompt:  "Project path",
					Default: "$HOME/project",
				},
			},
		},
	})

	require.NoError(t, err)

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	require.Len(t, valuesByRecipe, 1)
	assert.Equal(t, homeDir+"/project", valuesByRecipe[0]["project_path"])
}

func TestCollectRecipeInputsErrorsWithoutTTYAndDefault(t *testing.T) {
	_, err := collectRecipeInputs([]recipe.Recipe{
		{
			Name: "Input Recipe",
			Inputs: []recipe.Input{
				{
					ID:     "project_path",
					Kind:   recipe.InputPath,
					Prompt: "Project path",
				},
			},
		},
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "no TTY")
}

func TestCollectRecipeInputsScopesDuplicateIDsToRecipe(t *testing.T) {
	valuesByRecipe, err := collectRecipeInputs([]recipe.Recipe{
		{
			Name: "First",
			Inputs: []recipe.Input{
				{ID: "project_path", Kind: recipe.InputPath, Prompt: "First path", Default: "/first"},
			},
		},
		{
			Name: "Second",
			Inputs: []recipe.Input{
				{ID: "project_path", Kind: recipe.InputPath, Prompt: "Second path", Default: "/second"},
			},
		},
	})

	require.NoError(t, err)
	require.Len(t, valuesByRecipe, 2)
	assert.Equal(t, "/first", valuesByRecipe[0]["project_path"])
	assert.Equal(t, "/second", valuesByRecipe[1]["project_path"])
}

func TestCollectRecipeInputsSharesGlobalIDAcrossRecipes(t *testing.T) {
	valuesByRecipe, err := collectRecipeInputs([]recipe.Recipe{
		{
			Name: "First",
			Inputs: []recipe.Input{
				{ID: "project_root", GlobalID: "projects_path", Kind: recipe.InputPath, Prompt: "First path", Default: "/first"},
			},
		},
		{
			Name: "Second",
			Inputs: []recipe.Input{
				{ID: "workspace", GlobalID: "projects_path", Kind: recipe.InputPath, Prompt: "Second path", Default: "/second"},
			},
		},
	})

	require.NoError(t, err)
	require.Len(t, valuesByRecipe, 2)
	assert.Equal(t, "/first", valuesByRecipe[0]["project_root"])
	assert.Equal(t, "/first", valuesByRecipe[1]["workspace"])
}

func TestNormalizePathInputExpandsHome(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	assert.Equal(t, homeDir+"/project", normalizePathInput("~/project"))
}
