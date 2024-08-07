/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/yangyang5214/btl/pkg"

	"github.com/spf13/cobra"
)

// gpxSplitCmd represents the gpxSplit command
var gpxSplitCmd = &cobra.Command{
	Use:   "gpx_zip",
	Short: "zip gpx file to limited size.(remove some points)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gs := pkg.NewGpxZip(inputPath, step)
		err := gs.Run()
		if err != nil {
			panic(err)
		}
	},
}

var (
	step int
)

func init() {
	rootCmd.AddCommand(gpxSplitCmd)
	gpxSplitCmd.Flags().Int("s", step, "")
}
