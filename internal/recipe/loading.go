package recipe

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/floffah/culprit/internal/config"
	builtinRecipes "github.com/floffah/culprit/recipes"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

func LoadAllRecipes() ([]Recipe, error) {
	var recipes []Recipe

	embeddedEntries, err := builtinRecipes.RecipeFS.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range embeddedEntries {
		if entry.IsDir() {
			continue
		}

		data, err := builtinRecipes.RecipeFS.ReadFile(entry.Name())
		if err != nil {
			return nil, err
		}

		var recipe Recipe
		if err := toml.Unmarshal(data, &recipe); err != nil {
			return nil, err
		}

		recipe, err = ResolveInputPresets(recipe)
		if err != nil {
			return nil, err
		}

		ok, err := validateRecipe(recipe, entry.Name())
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		recipes = append(recipes, recipe)
	}

	// then read from viper config directory
	recipesDir := os.ExpandEnv(viper.GetString(config.RecipesKey))
	dirEntries, err := os.ReadDir(recipesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}

		source := filepath.Join(recipesDir, entry.Name())
		data, err := os.ReadFile(source)
		if err != nil {
			return nil, err
		}

		var recipe Recipe
		if err := toml.Unmarshal(data, &recipe); err != nil {
			return nil, err
		}

		recipe, err = ResolveInputPresets(recipe)
		if err != nil {
			return nil, err
		}

		ok, err := validateRecipe(recipe, source)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}
