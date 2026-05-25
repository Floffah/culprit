package recipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveInputPresetsPopulatesInputFields(t *testing.T) {
	resolved, err := ResolveInputPresets(Recipe{
		Name: "Preset Recipe",
		Inputs: []Input{
			{
				ID:     "workspace",
				Preset: "projects_path",
			},
		},
	})

	require.NoError(t, err)
	require.Len(t, resolved.Inputs, 1)

	input := resolved.Inputs[0]
	assert.Equal(t, "workspace", input.ID)
	assert.Equal(t, "projects_path", input.Preset)
	assert.Equal(t, "projects_path", input.GlobalID)
	assert.Equal(t, InputPath, input.Kind)
	assert.NotEmpty(t, input.Prompt)
	assert.NotEmpty(t, input.Default)
}

func TestResolveInputPresetsAllowsRecipeOverrides(t *testing.T) {
	resolved, err := ResolveInputPresets(Recipe{
		Name: "Preset Recipe",
		Inputs: []Input{
			{
				ID:       "workspace",
				Preset:   "projects_path",
				Prompt:   "Custom prompt",
				Default:  "~/Code",
				GlobalID: "custom_projects_path",
			},
		},
	})

	require.NoError(t, err)

	input := resolved.Inputs[0]
	assert.Equal(t, "custom_projects_path", input.GlobalID)
	assert.Equal(t, "Custom prompt", input.Prompt)
	assert.Equal(t, "~/Code", input.Default)
	assert.Equal(t, InputPath, input.Kind)
}

func TestResolveInputPresetsRejectsUnknownPreset(t *testing.T) {
	_, err := ResolveInputPresets(Recipe{
		Name: "Preset Recipe",
		Inputs: []Input{
			{
				ID:     "workspace",
				Preset: "missing",
			},
		},
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "unknown preset")
}
