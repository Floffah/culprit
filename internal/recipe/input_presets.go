package recipe

import "fmt"

var InputPresets = map[string]Input{
	"projects_path": {
		GlobalID:    "projects_path",
		Kind:        InputPath,
		Prompt:      "Where do you keep your projects?",
		Description: "Culprit will use this as the base path for project-specific cleanup recipes.",
		Default:     "~/Documents/Projects",
		Placeholder: "~/Documents/Projects",
	},
}

func ResolveInputPresets(recipe Recipe) (Recipe, error) {
	for inputIndex, input := range recipe.Inputs {
		if input.Preset == "" {
			continue
		}

		preset, ok := InputPresets[input.Preset]
		if !ok {
			return recipe, fmt.Errorf("recipe %q input %q references unknown preset %q", recipe.Name, input.ID, input.Preset)
		}

		recipe.Inputs[inputIndex] = mergeInputPreset(preset, input)
	}

	return recipe, nil
}

func mergeInputPreset(preset, input Input) Input {
	merged := preset

	merged.ID = input.ID
	merged.Preset = input.Preset

	if input.GlobalID != "" {
		merged.GlobalID = input.GlobalID
	}
	if input.Kind != "" {
		merged.Kind = input.Kind
	}
	if input.Prompt != "" {
		merged.Prompt = input.Prompt
	}
	if input.Description != "" {
		merged.Description = input.Description
	}
	if input.Default != "" {
		merged.Default = input.Default
	}
	if input.Placeholder != "" {
		merged.Placeholder = input.Placeholder
	}

	return merged
}
