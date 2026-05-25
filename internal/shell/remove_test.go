package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/floffah/culprit/internal/recipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemovesToBashCommentsEveryMetadataLine(t *testing.T) {
	script := RemovesToBash([]recipe.Remove{
		{
			RecipeName: "Unsafe\nprintf hacked",
			Reason:     "first line\nrm -rf /",
			Matcher:    "$HOME/cache\nopen /Applications",
			Path:       "/tmp/cache",
		},
	})

	for _, line := range strings.Split(script, "\n") {
		if strings.Contains(line, "printf hacked") && !strings.HasPrefix(line, "# ") {
			assert.Failf(t, "recipe metadata escaped comment", "%q", line)
		}
		if strings.Contains(line, "rm -rf /") && !strings.HasPrefix(line, "# ") {
			assert.Failf(t, "reason metadata escaped comment", "%q", line)
		}
		if strings.Contains(line, "open /Applications") && !strings.HasPrefix(line, "# ") {
			assert.Failf(t, "matcher metadata escaped comment", "%q", line)
		}
	}
}

func TestRemovesToBashShellQuotesPath(t *testing.T) {
	script := RemovesToBash([]recipe.Remove{
		{
			RecipeName: "Quoted Path",
			Reason:     "test",
			Matcher:    "test",
			Path:       "/tmp/it's a cache",
		},
	})

	want := "rm -rf -- '/tmp/it'\"'\"'s a cache'"
	assert.Contains(t, script, want)
}

func TestRemovesToBashWritesCommandAndSizesRelatedPaths(t *testing.T) {
	cacheDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(cacheDir, "one"), []byte("12345"), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(cacheDir, "nested"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(cacheDir, "nested", "two"), []byte("1234567"), 0644))

	script := RemovesToBash([]recipe.Remove{
		{
			RecipeName:   "Go Cache",
			Reason:       "Clean Go caches",
			Matcher:      "go clean -cache",
			Command:      "go clean -cache",
			RelatedPaths: []string{cacheDir},
		},
	})

	assert.Contains(t, script, "# Size: 12.0 B")
	assert.Contains(t, script, "\ngo clean -cache\n")
	assert.NotContains(t, script, "rm -rf --")
}

func TestRemovesToBashCanOmitSizes(t *testing.T) {
	script := RemovesToBashWithOptions([]recipe.Remove{
		{
			RecipeName: "Cache",
			Reason:     "test",
			Matcher:    "/tmp/cache",
			Path:       "/tmp/cache",
		},
	}, RenderOptions{IncludeSizes: false})

	assert.NotContains(t, script, "# Size:")
	assert.Contains(t, script, "rm -rf -- '/tmp/cache'")
}
