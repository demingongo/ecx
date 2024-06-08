/*
Copyright © 2024 demingongo
*/
package cmd

import (
	"github.com/demingongo/ecx/apps/applyapp"
	"github.com/demingongo/ecx/globals"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply ecx.yaml project file",
	Long: `Apply ecx.yaml project file.

The ecx.yaml project file should locate the resources to deploy.
For example:

ecx.yaml
+---------------------------------------------------+
| api: ecx                                          |
| apiVersion: 0.1                                   |
|                                                   |
| # cloudwatch log groups                           |
| logGroups:                                        |
|   - group: /etc/app-test                          |
|     retention: 1                                  |
|   - group: /etc/app2                              |
|                                                   |
| # ecs task definitions                            |
| taskDefinitions:                                  |
|   - taskdefinitions/taskdefinition.json           |
|   - taskdefinitions/taskdefinition2.json          |
|                                                   |
| # flows:                                          |
| #                                                 |
| # A flow could be                                 |
| # - rules, target group and service               |
| # - rules and target group                        |
| # - target group and service                      |
| # - target group                                  |
| # - or service                                    |
| #                                                 |
| # If you specify a target group for a service,    |
| # that service should have a container port       |
| # mapping named "http" in its task definition.    |
| flows:                                            |
|   - name: app-test-flow                           |
|     service: services/service.json                |
|     targetGroup: targetgroups/targetgroup.json    |
|     rules:                                        |
|       - value: rules/rule.json                    |
+---------------------------------------------------+
`,
	Run: func(cmd *cobra.Command, args []string) {
		globals.LoadGlobals()
		applyapp.Run()
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	applyCmd.PersistentFlags().StringP("project", "p", "", "path to the directory with ecx.yaml")
	applyCmd.MarkPersistentFlagDirname("project")

	viper.BindPFlag("project", applyCmd.PersistentFlags().Lookup("project"))
}
