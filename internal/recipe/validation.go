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

	seenInputIDs := make(map[string]struct{})
	for i, input := range recipe.Inputs {
		if strings.TrimSpace(input.ID) == "" {
			return false, fmt.Errorf("recipe %q input %d must define id", recipe.Name, i)
		}
		if _, exists := seenInputIDs[input.ID]; exists {
			return false, fmt.Errorf("recipe %q input %d duplicates id %q", recipe.Name, i, input.ID)
		}
		seenInputIDs[input.ID] = struct{}{}

		switch input.Kind {
		case InputPath:
			if strings.TrimSpace(input.Prompt) == "" {
				return false, fmt.Errorf("recipe %q input %q must define prompt", recipe.Name, input.ID)
			}
		default:
			return false, fmt.Errorf("recipe %q input %q has invalid kind %q", recipe.Name, input.ID, input.Kind)
		}
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
