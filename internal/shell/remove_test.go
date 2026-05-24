package shell

import (
	"strings"
	"testing"

	"github.com/floffah/maculprit/internal/recipe"
	"github.com/stretchr/testify/assert"
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
