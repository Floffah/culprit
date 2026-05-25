package clean

import (
	"context"
	"time"

	"charm.land/huh/v2/spinner"
	"github.com/floffah/culprit/internal/recipe"
	"github.com/floffah/culprit/internal/theming"
)

func discoverRecipes(ctx context.Context) ([]recipe.Recipe, error) {
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
			Context(ctx).
			Run()

		if err != nil {
			return nil, err
		}
	} else {
		loadedRecipes, err := recipe.LoadAllRecipes()
		if err != nil {
			return nil, err
		}
		recipes = loadedRecipes
	}

	return recipes, nil
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
