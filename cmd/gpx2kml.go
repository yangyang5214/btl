/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// gpx2kmlCmd represents the gpx2kml command
var gpx2kmlCmd = &cobra.Command{
	Use:   "gpx2kml",
	Short: "gpx to kml file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(gpx2kmlCmd)
}
