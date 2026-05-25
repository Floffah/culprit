package recipe

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type TemplateData struct {
	Inputs map[string]string
}

func ApplyInputValuesToRecipes(recipes []Recipe, valuesByRecipe []map[string]string) ([]Recipe, error) {
	if len(valuesByRecipe) == 0 {
		return recipes, nil
	}
	if len(valuesByRecipe) != len(recipes) {
		return nil, fmt.Errorf("got %d input value sets for %d recipes", len(valuesByRecipe), len(recipes))
	}

	resolved := make([]Recipe, 0, len(recipes))
	for recipeIndex, recipe := range recipes {
		resolvedRecipe, err := ApplyInputValues(recipe, valuesByRecipe[recipeIndex])
		if err != nil {
			return nil, err
		}
		resolved = append(resolved, resolvedRecipe)
	}
	return resolved, nil
}

func ApplyInputValues(recipe Recipe, values map[string]string) (Recipe, error) {
	data := TemplateData{Inputs: values}

	for targetIndex := range recipe.Targets {
		target := &recipe.Targets[targetIndex]

		path, err := renderRecipeTemplate(target.Path, data)
		if err != nil {
			return recipe, fmt.Errorf("recipe %q target %d path template: %w", recipe.Name, targetIndex, err)
		}
		target.Path = path

		command, err := renderRecipeTemplate(target.Command, data)
		if err != nil {
			return recipe, fmt.Errorf("recipe %q target %d command template: %w", recipe.Name, targetIndex, err)
		}
		target.Command = command

		for relatedPathIndex := range target.RelatedPaths {
			relatedPath, err := renderRecipeTemplate(target.RelatedPaths[relatedPathIndex], data)
			if err != nil {
				return recipe, fmt.Errorf("recipe %q target %d related_paths %d template: %w", recipe.Name, targetIndex, relatedPathIndex, err)
			}
			target.RelatedPaths[relatedPathIndex] = relatedPath
		}
	}

	return recipe, nil
}

func renderRecipeTemplate(value string, data TemplateData) (string, error) {
	if value == "" {
		return "", nil
	}

	tmpl, err := template.New("recipe-field").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"shellquote": shellQuote,
		}).
		Parse(value)
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", err
	}

	return rendered.String(), nil
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}

	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}
