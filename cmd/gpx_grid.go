/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/yangyang5214/btl/pkg/gpx_grid"

	"github.com/spf13/cobra"
)

// gpxGridCmd represents the gpxGrid command
var gpxGridCmd = &cobra.Command{
	Use:   "gpx_grid",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		grid := gpx_grid.NewGpxGrid()
		err := grid.Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gpxGridCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gpxGridCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gpxGridCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
