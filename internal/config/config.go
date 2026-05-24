package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

const RecipesKey = "recipes"

func SetupConfig() error {
	viper.SetConfigName("maculprit")
	viper.SetEnvPrefix("maculprit")
	viper.AutomaticEnv()

	viper.AddConfigPath("$HOME/.maculprit")
	viper.AddConfigPath(".")

	viper.SetConfigType("toml")

	viper.SetDefault(RecipesKey, "$HOME/.maculprit/recipes")

	err := os.MkdirAll(os.ExpandEnv("$HOME/.maculprit"), 0755)
	if err != nil {
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = os.MkdirAll(os.ExpandEnv(viper.GetString(RecipesKey)), 0755)
	if err != nil {
		return err
	}

	return nil
}
