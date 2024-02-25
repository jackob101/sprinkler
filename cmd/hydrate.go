/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	hydrator "jackob101/sprinkler/lib"

	"github.com/spf13/cobra"
)

// hydrateCmd represents the hydrate command
var hydrateCmd = &cobra.Command{
	Use:   "hydrate [path to file with variables] [path to directory with templates]",
	Short: "Fires off hydration",
	Long:  `TODO add longer description`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			panic("Failed to get output flag")
		}

		println(outputPath)
		if outputPath == "" {
			outputPath = "filled_templates"
		}

		pathToVariables := args[0]
		pathToTemplates := args[1]

		hydrator.Hydrate(pathToVariables, pathToTemplates, outputPath)
	},
}

func init() {
	rootCmd.AddCommand(hydrateCmd)

	hydrateCmd.Flags().String("output", "", "Output directory")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hydrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hydrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
