package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

const RecipesKey = "recipes"

func SetupConfig() error {
	viper.SetConfigName("culprit")
	viper.SetEnvPrefix("culprit")
	viper.AutomaticEnv()

	viper.AddConfigPath("$HOME/.culprit")
	viper.AddConfigPath(".")

	viper.SetConfigType("toml")

	viper.SetDefault(RecipesKey, "$HOME/.culprit/recipes")

	err := os.MkdirAll(os.ExpandEnv("$HOME/.culprit"), 0755)
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
