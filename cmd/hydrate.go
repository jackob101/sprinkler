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
	Use:   "hydrate",
	Short: "Fires off hydration",
	Long:  `TODO add longer description`,
	Run: func(cmd *cobra.Command, args []string) {
		hydrator.Hydrate()
	},
}

func init() {
	rootCmd.AddCommand(hydrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hydrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hydrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
