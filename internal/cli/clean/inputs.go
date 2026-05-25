package clean

import (
	"fmt"
	"os"
	"strings"

	"charm.land/huh/v2"
	"github.com/floffah/culprit/internal/recipe"
	"github.com/floffah/culprit/internal/theming"
)

func collectRecipeInputs(recipes []recipe.Recipe) ([]map[string]string, error) {
	isTTY := theming.AreWeTTY()

	valuesByRecipe := make([]map[string]string, len(recipes))
	globalValues := make(map[string]string)
	globalInputs := make(map[string]recipe.Input)

	for recipeIndex, r := range recipes {
		valuesByRecipe[recipeIndex] = make(map[string]string, len(r.Inputs))

		for _, input := range r.Inputs {
			value, err := collectScopedInputValue(input, globalValues, globalInputs, isTTY)
			if err != nil {
				return nil, err
			}
			valuesByRecipe[recipeIndex][input.ID] = value
		}
	}

	return valuesByRecipe, nil
}

func collectScopedInputValue(input recipe.Input, globalValues map[string]string, globalInputs map[string]recipe.Input, isTTY bool) (string, error) {
	if input.GlobalID == "" {
		return collectInputValue(input, isTTY)
	}

	if existingInput, ok := globalInputs[input.GlobalID]; ok {
		if existingInput.Kind != input.Kind {
			return "", fmt.Errorf("global input %q has conflicting kinds %q and %q", input.GlobalID, existingInput.Kind, input.Kind)
		}
		return globalValues[input.GlobalID], nil
	}

	value, err := collectInputValue(input, isTTY)
	if err != nil {
		return "", err
	}

	globalInputs[input.GlobalID] = input
	globalValues[input.GlobalID] = value
	return value, nil
}

func collectInputValue(input recipe.Input, isTTY bool) (string, error) {
	switch input.Kind {
	case recipe.InputPath:
		return collectPathInput(input, isTTY)
	default:
		return "", fmt.Errorf("input %q has unsupported kind %q", input.ID, input.Kind)
	}
}

func collectPathInput(input recipe.Input, isTTY bool) (string, error) {
	value := input.Default

	if !isTTY {
		if strings.TrimSpace(value) == "" {
			return "", fmt.Errorf("input %q requires a path but no TTY is available and no default was provided", input.ID)
		}
		return normalizePathInput(value), nil
	}

	field := huh.NewInput().
		Title(input.Prompt).
		Description(input.Description).
		Placeholder(input.Placeholder).
		Value(&value).
		Validate(func(value string) error {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("path is required")
			}
			return nil
		})

	if err := field.Run(); err != nil {
		return "", err
	}

	return normalizePathInput(value), nil
}

func normalizePathInput(value string) string {
	value = strings.TrimSpace(os.ExpandEnv(value))
	if value == "~" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			return homeDir
		}
	}

	if strings.HasPrefix(value, "~/") {
		if homeDir, err := os.UserHomeDir(); err == nil {
			return homeDir + value[1:]
		}
	}

	return value
}
