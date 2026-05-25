package recipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeriveRecipeRemovesCreatesCommandRemove(t *testing.T) {
	removes, err := DeriveRecipeRemoves(Recipe{
		Name: "Go Cache",
		Type: TypeSoft,
		Targets: []Target{
			{
				Reason:       "Clean Go caches",
				Kind:         TargetCommand,
				Command:      "go clean -cache && go clean -modcache",
				RelatedPaths: []string{"/tmp/go-cache"},
			},
		},
	})

	require.NoError(t, err)
	require.Len(t, removes, 1)
	assert.Equal(t, "Go Cache", removes[0].RecipeName)
	assert.Equal(t, "go clean -cache && go clean -modcache", removes[0].Command)
	assert.Equal(t, "go clean -cache && go clean -modcache", removes[0].Matcher)
	assert.Equal(t, []string{"/tmp/go-cache"}, removes[0].RelatedPaths)
	assert.Empty(t, removes[0].Path)
}

func TestDeriveRecipeRemovesRejectsCommandWithRM(t *testing.T) {
	_, err := DeriveRecipeRemoves(Recipe{
		Name: "Unsafe Command",
		Type: TypeSoft,
		Targets: []Target{
			{
				Reason:       "test",
				Kind:         TargetCommand,
				Command:      "sudo /bin/rm -rf /tmp/cache",
				RelatedPaths: []string{"/tmp/cache"},
			},
		},
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "includes rm")
}
