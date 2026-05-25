package clean

import (
	"fmt"

	"charm.land/log/v2"
	"github.com/floffah/culprit/internal/recipe"
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
	cmd.Flags().Bool("no-sizes", false, "Skip calculating clearable sizes while generating the cleanup script (can speed up generation significantly)")

	return cmd
}

func RunE(cmd *cobra.Command, args []string) error {
	recipes, err := discoverRecipes(cmd.Context())
	if err != nil {
		return err
	}
	log.Infof("Discovered %d recipes", len(recipes))

	givenType := recipeTypeFromArgs(args)
	IsCleaningCulprit = givenType == string(recipe.TypeCulprit)

	usableRecipes := filterRecipesByType(recipes, givenType)

	log.Infof("%d recipes are of type %s", len(usableRecipes), givenType)

	inputValues, err := collectRecipeInputs(usableRecipes)
	if err != nil {
		return err
	}
	usableRecipes, err = recipe.ApplyInputValuesToRecipes(usableRecipes, inputValues)
	if err != nil {
		return err
	}

	allRemoves, err := computeRemoves(usableRecipes, givenType)
	if err != nil {
		return err
	}

	log.Infof("Derived %d removes from %d recipes", len(allRemoves), len(usableRecipes))

	noSizes, err := cmd.Flags().GetBool("no-sizes")
	if err != nil {
		return err
	}

	script, err := buildScript(allRemoves, cmd.Context(), !noSizes)
	if err != nil {
		return err
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

	isTTY := theming.AreWeTTY()
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
