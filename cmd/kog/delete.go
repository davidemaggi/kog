/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package kog

import (
	"github.com/davidemaggi/kog/pkg/wizard"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var deleteCmd = &cobra.Command{
	Version: version,
	Use:     "delete",
	Aliases: []string{"d"},
	Short:   "Delete the selected context",
	Long: `This command deletes the selected context from your kubeconfig file.
Clusters, Users etc... remain untouched`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		wizard.DeleteContext(ConfigPath, Verbose)

	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wizard.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(deleteCmd)

}
