/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/yangyang5214/btl/pkg/marker_point"

	"github.com/spf13/cobra"
)

// markerPointCmd represents the markerPoint command
var markerPointCmd = &cobra.Command{
	Use:   "markerPoint",
	Short: "amap marker point",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		mp := marker_point.NewMarkerPoint(pointsFile)
		err := mp.Run()
		if err != nil {
			panic(err)
		}
	},
}

var pointsFile string

func init() {
	rootCmd.AddCommand(markerPointCmd)
	markerPointCmd.Flags().StringVarP(&pointsFile, "point_file", "f", "", "points file")
}
