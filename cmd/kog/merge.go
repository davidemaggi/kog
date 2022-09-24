/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package kog

import (
	"github.com/davidemaggi/kog/pkg/merge"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var mergeCmd = &cobra.Command{
	Version: version,
	Use:     "merge",
	Aliases: []string{"m"},
	Args:    cobra.ExactArgs(1),
	Short:   "Merge a new Yaml config to your existing KubeConfig",
	Long: `This is the command that i was missing from other tools...
Switching from multiple environments was a pain in the ass...`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		//log.Print(args[0])
		merge.MergeConfig(args[0], ConfigPath, Force, Verbose)
	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wizard.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(mergeCmd)

}
