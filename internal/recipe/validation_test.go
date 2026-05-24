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
