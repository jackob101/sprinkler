/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"jackob101/sprinkler/lib"

	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describes what will be changed in templates",
	Long:  `Returns all variables from templates and values that will replace them. Used to check if there are no missing variables`,
	Run: func(cmd *cobra.Command, args []string) {
		pathToVariables := args[0]
		pathToTemplates := args[1]

		lib.Describe(pathToVariables, pathToTemplates)
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// describeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// describeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
