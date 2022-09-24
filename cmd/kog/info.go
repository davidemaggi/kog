/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package kog

import (
	"github.com/davidemaggi/kog/pkg/wizard"
	"github.com/spf13/cobra"
)

var raw bool

// rootCmd represents the base command when called without any subcommands
var infoCmd = &cobra.Command{
	Version: version,
	Use:     "info",
	Aliases: []string{"i", "current"},
	Short:   "Get the Current Configuration",
	Long:    `Display the current Context and namespace.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		wizard.ShowInfo(ConfigPath, Verbose, raw)

	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wizard.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	infoCmd.Flags().BoolVar(&raw, "raw", false, "Show the entire config")
	rootCmd.AddCommand(infoCmd)

}
