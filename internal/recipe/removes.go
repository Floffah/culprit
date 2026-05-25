package recipe

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Remove struct {
	RecipeName   string
	Reason       string
	Matcher      string
	Path         string
	Command      string
	RelatedPaths []string
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
		case TargetCommand:
			if commandIncludesRM(target.Command) {
				return nil, fmt.Errorf("recipe %q command target includes rm", recipe.Name)
			}

			removes = append(removes, Remove{
				RecipeName:   recipe.Name,
				Reason:       target.Reason,
				Matcher:      target.Command,
				Command:      target.Command,
				RelatedPaths: target.RelatedPaths,
			})
		default:
			return nil, fmt.Errorf("recipe %q has unsupported target kind %q", recipe.Name, target.Kind)
		}
	}

	return removes, nil
}

var rmCommandPattern = regexp.MustCompile(`(^|[;&|()\s])(?:sudo\s+)?(?:/[^;&|()\s]+/)?rm([;&|()\s]|$)`)

func commandIncludesRM(command string) bool {
	return rmCommandPattern.MatchString(command)
}
