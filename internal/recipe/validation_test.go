package recipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRecipeSkipsInvalidType(t *testing.T) {
	ok, err := validateRecipe(Recipe{
		Name: "Bad Type",
		Type: Type("weird"),
		Targets: []Target{
			{
				Reason: "test",
				Kind:   TargetAbsolute,
				Path:   "/tmp/cache",
			},
		},
	}, "bad.recipe.toml")

	require.NoError(t, err)
	assert.False(t, ok)
}

func TestValidateRecipeRejectsInvalidTargetKind(t *testing.T) {
	ok, err := validateRecipe(Recipe{
		Name: "Bad Target",
		Type: TypeSoft,
		Targets: []Target{
			{
				Reason: "test",
				Kind:   TargetKind("mystery"),
				Path:   "/tmp/cache",
			},
		},
	}, "bad.recipe.toml")

	require.Error(t, err)
	assert.False(t, ok)
}

func TestValidateRecipeAcceptsCommandTarget(t *testing.T) {
	ok, err := validateRecipe(Recipe{
		Name: "Command Target",
		Type: TypeSoft,
		Targets: []Target{
			{
				Reason:       "test",
				Kind:         TargetCommand,
				Command:      "go clean -cache && terraform version",
				RelatedPaths: []string{"/tmp/cache"},
			},
		},
	}, "command.recipe.toml")

	require.NoError(t, err)
	assert.True(t, ok)
}

func TestValidateRecipeRejectsCommandTargetWithRM(t *testing.T) {
	ok, err := validateRecipe(Recipe{
		Name: "Command Target",
		Type: TypeSoft,
		Targets: []Target{
			{
				Reason:       "test",
				Kind:         TargetCommand,
				Command:      "go clean -cache && rm -rf /tmp/cache",
				RelatedPaths: []string{"/tmp/cache"},
			},
		},
	}, "command.recipe.toml")

	require.Error(t, err)
	assert.False(t, ok)
}

func TestValidateRecipeAcceptsPathInput(t *testing.T) {
	ok, err := validateRecipe(Recipe{
		Name: "Path Input",
		Type: TypeSoft,
		Inputs: []Input{
			{
				ID:     "project_path",
				Kind:   InputPath,
				Prompt: "Project path",
			},
		},
		Targets: []Target{
			{
				Reason: "test",
				Kind:   TargetAbsolute,
				Path:   "{{ .Inputs.project_path }}/cache",
			},
		},
	}, "path-input.recipe.toml")

	require.NoError(t, err)
	assert.True(t, ok)
}

func TestValidateRecipeRejectsInvalidInputKind(t *testing.T) {
	ok, err := validateRecipe(Recipe{
		Name: "Path Input",
		Type: TypeSoft,
		Inputs: []Input{
			{
				ID:     "project_path",
				Kind:   InputKind("mystery"),
				Prompt: "Project path",
			},
		},
		Targets: []Target{
			{
				Reason: "test",
				Kind:   TargetAbsolute,
				Path:   "{{ .Inputs.project_path }}/cache",
			},
		},
	}, "path-input.recipe.toml")

	require.Error(t, err)
	assert.False(t, ok)
}
