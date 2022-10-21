package kog

import (
	"fmt"
	"os"

	wizard "github.com/davidemaggi/kog/pkg/wizard"
	"github.com/spf13/cobra"
)

var version = "1.1.1"
var Verbose bool
var Force bool
var ConfigPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: version,
	Use:     "kog",
	Short:   "Yet another Kubernetes configuration manager.",
	Long:    `From a lazy developer to other lazy developers, a quick way to manage your kube config file adding, removing switching environment in a lazy way`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		wizard.SelectContext(ConfigPath, Verbose)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {

	version = v

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wizard.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "", "KubeConfig Path")
	rootCmd.PersistentFlags().BoolVarP(&Force, "force", "f", false, "Force Flag")

}
