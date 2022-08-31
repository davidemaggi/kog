/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package kog

import (
	"log"

	"github.com/davidemaggi/kog/pkg/k8s"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var mergeCmd = &cobra.Command{
	Version: version,
	Use:     "merge",
	Aliases: []string{"m"},
	Args:    cobra.ExactArgs(1),
	Short:   "Merge a new Yaml to your existing KubeConfig",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		log.Print(args[0])
		xxx, _ := k8s.FindKubeConfig()
		k8s.MergeConfigs(args[0], xxx, false, Verbose)
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
