/*
Copyright © 2024 demingongo
*/
package cmd

import (
	"github.com/demingongo/ecx/apps/updateserviceapp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateServiceCmd represents the updateService command
var updateServiceCmd = &cobra.Command{
	Use:   "update-service",
	Short: "Updates an ECS service",
	Long: `The command updates an ECS service.

It helps you:
	selecting new tags for containers in the task(s),
	creating new revisions of the task definition(s),
	updating the service with the new revisions.`,
	Run: func(cmd *cobra.Command, args []string) {
		updateserviceapp.Run()
	},
}

func init() {
	rootCmd.AddCommand(updateServiceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateServiceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateServiceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateServiceCmd.PersistentFlags().String("cluster", "", "cluster name")
	updateServiceCmd.PersistentFlags().String("service", "", "ecs service arn")
	updateServiceCmd.MarkFlagsMutuallyExclusive("cluster", "service")

	viper.BindPFlag("cluster", updateServiceCmd.PersistentFlags().Lookup("cluster"))
	viper.BindPFlag("service", updateServiceCmd.PersistentFlags().Lookup("service"))
}
