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
	"github.com/spf13/cobra"
)

var IsCleaningCulprit bool

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean [type]",
		Short: "Evaluate all recipes and creates a shell script to clean up derived files",
		Args:  cobra.MatchAll(cobra.OnlyValidArgs),
		RunE:  RunE,
	}
	cmd.Flags().StringP("output", "o", "cleanup.sh", "The output file to write the cleanup script to")
	cmd.Flags().Bool("force", false, "Overwrite the output file if it already exists")

	return cmd
}

func RunE(cmd *cobra.Command, args []string) error {
	var recipes []recipe.Recipe
	isTTY := theming.AreWeTTY()

	if isTTY {
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
	} else {
		loadedRecipes, err := recipe.LoadAllRecipes()
		if err != nil {
			return err
		}
		recipes = loadedRecipes
	}

	log.Infof("Discovered %d recipes", len(recipes))

	givenType := recipeTypeFromArgs(args)
	IsCleaningCulprit = givenType == string(recipe.TypeCulprit)

	usableRecipes := filterRecipesByType(recipes, givenType)

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

	if isTTY {
		err := spinner.New().
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
	} else {
		script = shell.RemovesToBash(allRemoves)
	}

	outputFilePath, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}

	if err := writeScriptFile(outputFilePath, script, force); err != nil {
		return err
	}

	if isTTY {
		fmt.Println()
		fmt.Println(theming.ScriptWrittenStyle().Render("Cleanup script created!"))
		fmt.Print(theming.ScriptRunScriptBeforeStyle().Render("Run `"))
		fmt.Print(theming.ScriptRunCommandStyle().Render("bash " + outputFilePath))
		fmt.Println(theming.ScriptRunScriptAfterStyle().Render("` to execute the cleanup script"))
		fmt.Println()
		fmt.Println(theming.DangerStyle().Render("⚠️  Make sure to review the cleanup script before running it, as it may contain destructive commands!"))
	} else {
		log.Infof("Cleanup script written to %q", outputFilePath)
		log.Warnf("Make sure to review the cleanup script before running it, as it may contain destructive commands!")
	}

	return nil
}

func recipeTypeFromArgs(args []string) string {
	if len(args) == 0 {
		return string(recipe.TypeSoft)
	}

	return args[0]
}

func filterRecipesByType(recipes []recipe.Recipe, recipeType string) []recipe.Recipe {
	var filtered []recipe.Recipe
	for _, r := range recipes {
		if string(r.Type) == recipeType {
			filtered = append(filtered, r)
		}
	}
	return filtered
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
