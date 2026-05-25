package clean

import (
	"context"
	"fmt"
	"time"

	"charm.land/huh/v2/spinner"
	"github.com/floffah/culprit/internal/recipe"
	"github.com/floffah/culprit/internal/theming"
)

func computeRemoves(recipes []recipe.Recipe, givenType string) ([]recipe.Remove, error) {
	isTTY := theming.AreWeTTY()

	var allRemoves []recipe.Remove

	if isTTY {
		err := spinner.New().
			Title("Discovering removable items...").
			WithTheme(theming.SpinnerTheme()).
			ActionWithErr(func(ctx context.Context) error {
				startedAt := time.Now()

				derivedRemoves, err := inner(recipes, givenType)
				if err != nil {
					return err
				}
				allRemoves = derivedRemoves

				if time.Since(startedAt) < 500*time.Millisecond {
					time.Sleep(500*time.Millisecond - time.Since(startedAt))
				}

				return nil
			}).
			Context(context.Background()).
			Run()

		if err != nil {
			return nil, err
		}
	} else {
		var err error
		allRemoves, err = inner(recipes, givenType)
		if err != nil {
			return nil, err
		}
	}

	return allRemoves, nil
}

func inner(recipes []recipe.Recipe, givenType string) ([]recipe.Remove, error) {
	var allRemoves []recipe.Remove

	for _, r := range recipes {
		removes, err := recipe.DeriveRecipeRemoves(r)
		if err != nil {
			return nil, err
		}
		allRemoves = append(allRemoves, removes...)
	}

	if len(allRemoves) == 0 {
		return nil, fmt.Errorf("no removes derived from %d recipes of type %s", len(recipes), givenType)
	}

	return allRemoves, nil
}
