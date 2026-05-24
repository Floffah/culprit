package recipe

import (
	"fmt"
	"os"
	"path/filepath"
)

type Remove struct {
	RecipeName string
	Reason     string
	Matcher    string
	Path       string
}

func DeriveRecipeRemoves(recipe Recipe) ([]Remove, error) {
	var removes []Remove

	for _, target := range recipe.Targets {
		switch target.Kind {
		case TargetAbsolute:
			matches, err := filepath.Glob(os.ExpandEnv(target.Path))
			if err != nil {
				return nil, err
			}

			for _, match := range matches {
				removes = append(removes, Remove{
					RecipeName: recipe.Name,
					Reason:     target.Reason,
					Matcher:    target.Path,
					Path:       match,
				})
			}
		default:
			return nil, fmt.Errorf("recipe %q has unsupported target kind %q", recipe.Name, target.Kind)
		}
	}

	return removes, nil
}
