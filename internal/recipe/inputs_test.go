package recipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyInputValuesReplacesTargetTokens(t *testing.T) {
	resolved, err := ApplyInputValues(Recipe{
		Name: "Input Recipe",
		Targets: []Target{
			{
				Kind:         TargetCommand,
				Path:         "{{ .Inputs.project_path }}/cache",
				Command:      "tool clean {{ shellquote .Inputs.project_path }}",
				RelatedPaths: []string{"{{ .Inputs.project_path }}/build", "{{ .Inputs.project_path }}/cache"},
			},
		},
	}, map[string]string{
		"project_path": "/tmp/project",
	})
	require.NoError(t, err)

	target := resolved.Targets[0]
	assert.Equal(t, "/tmp/project/cache", target.Path)
	assert.Equal(t, "tool clean '/tmp/project'", target.Command)
	assert.Equal(t, []string{"/tmp/project/build", "/tmp/project/cache"}, target.RelatedPaths)
}

func TestApplyInputValuesShellQuotesCommandTokens(t *testing.T) {
	resolved, err := ApplyInputValues(Recipe{
		Name: "Input Recipe",
		Targets: []Target{
			{
				Kind:    TargetCommand,
				Command: "tool clean {{ shellquote .Inputs.project_path }}",
			},
		},
	}, map[string]string{
		"project_path": "/tmp/it's complicated",
	})
	require.NoError(t, err)

	assert.Equal(t, "tool clean '/tmp/it'\"'\"'s complicated'", resolved.Targets[0].Command)
}

func TestApplyInputValuesErrorsForMissingInput(t *testing.T) {
	_, err := ApplyInputValues(Recipe{
		Name: "Input Recipe",
		Targets: []Target{
			{
				Kind: TargetAbsolute,
				Path: "{{ .Inputs.missing }}/cache",
			},
		},
	}, map[string]string{
		"project_path": "/tmp/project",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "map has no entry for key")
}

func TestApplyInputValuesErrorsForInvalidTemplate(t *testing.T) {
	_, err := ApplyInputValues(Recipe{
		Name: "Input Recipe",
		Targets: []Target{
			{
				Kind: TargetAbsolute,
				Path: "{{ .Inputs.project_path ",
			},
		},
	}, map[string]string{
		"project_path": "/tmp/project",
	})

	require.Error(t, err)
}

func TestApplyInputValuesToRecipesUsesRecipeScopedValues(t *testing.T) {
	resolved, err := ApplyInputValuesToRecipes([]Recipe{
		{
			Name: "First",
			Targets: []Target{
				{Kind: TargetAbsolute, Path: "{{ .Inputs.project_path }}/cache"},
			},
		},
		{
			Name: "Second",
			Targets: []Target{
				{Kind: TargetAbsolute, Path: "{{ .Inputs.project_path }}/cache"},
			},
		},
	}, []map[string]string{
		{"project_path": "/first"},
		{"project_path": "/second"},
	})

	require.NoError(t, err)
	require.Len(t, resolved, 2)
	assert.Equal(t, "/first/cache", resolved[0].Targets[0].Path)
	assert.Equal(t, "/second/cache", resolved[1].Targets[0].Path)
}
