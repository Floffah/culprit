package clean

import (
	"context"
	"fmt"
	"os"
	"time"

	"charm.land/huh/v2/spinner"
	"charm.land/log/v2"
	"github.com/floffah/culprit/internal/recipe"
	"github.com/floffah/culprit/internal/shell"
	"github.com/floffah/culprit/internal/theming"
)

func buildScript(removes []recipe.Remove, ctx context.Context, includeSizes bool) (string, error) {
	isTTY := theming.AreWeTTY()

	var script string

	if isTTY {
		err := spinner.New().
			Title("Generating cleanup script...").
			WithTheme(theming.SpinnerTheme()).
			ActionWithErr(func(ctx context.Context) error {
				startedAt := time.Now()

				script = shell.RemovesToBashWithOptions(removes, shell.RenderOptions{
					IncludeSizes: includeSizes,
				})

				if time.Since(startedAt) < 500*time.Millisecond {
					time.Sleep(500*time.Millisecond - time.Since(startedAt))
				}

				return nil
			}).
			Context(ctx).
			Run()

		if err != nil {
			return "", err
		}
	} else {
		script = shell.RemovesToBashWithOptions(removes, shell.RenderOptions{
			IncludeSizes: includeSizes,
		})
	}

	return script, nil
}

func writeScriptFile(path, script string, force bool) error {
	openFlags := os.O_CREATE | os.O_WRONLY | os.O_EXCL
	if force {
		openFlags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	outputFile, err := os.OpenFile(path, openFlags, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("output file %q already exists, pass --force to overwrite it", path)
		}
		return err
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			log.Errorf("Failed to close output file: %v", err)
		}
	}(outputFile)

	_, err = outputFile.WriteString(script)
	return err
}
