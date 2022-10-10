/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package kog

import (
	"github.com/davidemaggi/kog/pkg/wizard"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var fwdCmd = &cobra.Command{
	Version: version,
	Use:     "forward",
	Aliases: []string{"f", "fwd"},
	Short:   "Forward a service on localhost",
	Long:    `A wizard to guide you on port forwarding configuration`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		err := wizard.PortForwarding(ConfigPath, Verbose)
		if err != nil {
			return
		}

	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wizard.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	//fwdCmd.Flags().BoolVar(&raw, "raw", false, "Show the entire config")
	rootCmd.AddCommand(fwdCmd)

}
