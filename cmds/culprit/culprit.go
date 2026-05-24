package main

import (
	"context"
	"os"

	"charm.land/fang/v2"
	"charm.land/log/v2"
	"github.com/floffah/culprit/internal/cli/clean"
	"github.com/floffah/culprit/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isVerbose bool

var rootCmd = &cobra.Command{
	Use: "culprit",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if isVerbose {
			log.SetLevel(log.DebugLevel)
		}

		err := config.SetupConfig()
		if err != nil {
			return err
		}

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if !clean.IsCleaningCulprit {
			return viper.WriteConfig()
		}

		return nil
	},
}

func main() {
	log.SetReportCaller(false)
	log.SetReportTimestamp(false)

	cobra.EnableTraverseRunHooks = true

	rootCmd.PersistentFlags().BoolVar(&isVerbose, "verbose", false, "Enable verbose theming")
	rootCmd.AddCommand(clean.NewCommand())

	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}
