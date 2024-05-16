/*
Copyright © 2024 demingongo
*/
package cmd

import (
	"os"

	"github.com/demingongo/ecx/apps/starterapp"
	"github.com/demingongo/ecx/globals"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ecx",
	Short: "A helper for AWS EC2/ECR/ECS",
	Long: `
    .----.
    |~_  |
  __|____|__
 |  ______--|
 '-/.::::.\-'
  '--------'
Using AWS EC2/ECR/ECS?
Tired of using the slow AWS console?

Here's a helper that uses aws-cli under the hood.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		globals.LoadGlobals()
		starterapp.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ecx.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("dummy", false, "dummy run (no aws call)")
	rootCmd.PersistentFlags().BoolP("colors", "c", false, "colorful forms")

	viper.BindPFlag("dummy", rootCmd.PersistentFlags().Lookup("dummy"))
	viper.BindPFlag("colors", rootCmd.PersistentFlags().Lookup("colors"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.SetDefault("dummy", false)
	viper.SetDefault("verbose", false)
}
