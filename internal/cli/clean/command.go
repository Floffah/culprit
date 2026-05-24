package clean

import (
	"context"
	"fmt"
	"os"
	"time"

	"charm.land/huh/v2/spinner"
	"charm.land/log/v2"
	"github.com/floffah/maculprit/internal/recipe"
	"github.com/floffah/maculprit/internal/shell"
	"github.com/floffah/maculprit/internal/theming"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean [type]",
		Short: "Evaluate all recipes and creates a shell script to clean up derived files",
		Args:  cobra.MatchAll(cobra.OnlyValidArgs),
		RunE:  RunE,
	}
	cmd.Flags().String("output", "cleanup.sh", "The output file to write the cleanup script to")

	return cmd
}

func RunE(cmd *cobra.Command, args []string) error {
	var recipes []recipe.Recipe

	err := spinner.New().
		Title("Discovering recipes...").
		WithTheme(theming.SpinnerTheme()).
		ActionWithErr(func(ctx context.Context) error {
			startedAt := time.Now()

			loadedRecipes, err := recipe.LoadAllRecipes()
			if err != nil {
				return err
			}
			recipes = loadedRecipes

			if time.Since(startedAt) < 500*time.Millisecond {
				time.Sleep(500*time.Millisecond - time.Since(startedAt))
			}

			return nil
		}).
		Context(cmd.Context()).
		Run()

	if err != nil {
		return err
	}

	log.Infof("Discovered %d recipes", len(recipes))

	givenType := "soft"

	if len(args) > 0 {
		givenType = args[0]
	}

	var usableRecipes []recipe.Recipe
	for _, r := range recipes {
		if string(r.Type) == givenType {
			usableRecipes = append(usableRecipes, r)
		}
	}

	log.Infof("%d recipes are of type %s", len(usableRecipes), givenType)

	var allRemoves []recipe.Remove
	for _, r := range usableRecipes {
		removes, err := recipe.DeriveRecipeRemoves(r)
		if err != nil {
			return err
		}
		allRemoves = append(allRemoves, removes...)
	}

	if len(allRemoves) == 0 {
		log.Infof("No removes derived from recipes of type %s", givenType)
		return nil
	}

	log.Infof("Derived %d removes from %d recipes", len(allRemoves), len(usableRecipes))

	var script string

	err = spinner.New().
		Title("Generating cleanup script...").
		WithTheme(theming.SpinnerTheme()).
		ActionWithErr(func(ctx context.Context) error {
			startedAt := time.Now()

			script = shell.RemovesToBash(allRemoves)

			if time.Since(startedAt) < 500*time.Millisecond {
				time.Sleep(500*time.Millisecond - time.Since(startedAt))
			}

			return nil
		}).
		Context(cmd.Context()).
		Run()

	if err != nil {
		return err
	}

	outputFilePath, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(script)
	if err != nil {
		return err
	}

	printCleanResult(outputFilePath)

	return nil
}

func printCleanResult(atPath string) {
	fmt.Println()
	fmt.Println(theming.MutedStyle().Render("⎯⎯⎯⎯⎯"))
	fmt.Println()
	fmt.Println(theming.ScriptWrittenStyle().Render("Cleanup script created!"))
	fmt.Print(theming.ScriptRunScriptBeforeStyle().Render("Run `"))
	fmt.Print(theming.ScriptRunCommandStyle().Render("bash " + atPath))
	fmt.Println(theming.ScriptRunScriptAfterStyle().Render("` to execute the cleanup script"))
	fmt.Println()
	fmt.Println(theming.DangerStyle().Render("⚠️  Make sure to review the cleanup script before running it, as it may contain destructive commands!"))
}
