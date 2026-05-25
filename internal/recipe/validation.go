package recipe

import (
	"fmt"
	"strings"

	"charm.land/log/v2"
)

func validateRecipe(recipe Recipe, source string) (bool, error) {
	displayName := recipe.Name
	if displayName == "" {
		displayName = source
	}

	if !validType(recipe.Type) {
		log.Warnf("Skipping recipe %q from %s: invalid type %q", displayName, source, recipe.Type)
		return false, nil
	}

	if strings.TrimSpace(recipe.Name) == "" {
		return false, fmt.Errorf("recipe from %s must define name", source)
	}

	if len(recipe.Targets) == 0 {
		return false, fmt.Errorf("recipe %q from %s must define at least one target", recipe.Name, source)
	}

	for i, target := range recipe.Targets {
		if strings.TrimSpace(target.Reason) == "" {
			return false, fmt.Errorf("recipe %q target %d must define reason", recipe.Name, i)
		}

		switch target.Kind {
		case TargetAbsolute:
			if strings.TrimSpace(target.Path) == "" {
				return false, fmt.Errorf("recipe %q target %d must define path", recipe.Name, i)
			}
		case TargetCommand:
			if strings.TrimSpace(target.Command) == "" {
				return false, fmt.Errorf("recipe %q target %d must define command", recipe.Name, i)
			}
			if commandIncludesRM(target.Command) {
				return false, fmt.Errorf("recipe %q target %d command must not include rm", recipe.Name, i)
			}
			if len(target.RelatedPaths) == 0 {
				return false, fmt.Errorf("recipe %q target %d must define related_paths", recipe.Name, i)
			}
			for relatedPathIndex, relatedPath := range target.RelatedPaths {
				if strings.TrimSpace(relatedPath) == "" {
					return false, fmt.Errorf("recipe %q target %d related_paths entry %d must not be empty", recipe.Name, i, relatedPathIndex)
				}
			}
		default:
			return false, fmt.Errorf("recipe %q target %d has invalid kind %q", recipe.Name, i, target.Kind)
		}
	}

	return true, nil
}

func validType(recipeType Type) bool {
	switch recipeType {
	case TypeSoft, TypeHard, TypeCulprit:
		return true
	default:
		return false
	}
}
