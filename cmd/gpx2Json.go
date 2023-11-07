/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

var gpx2JsonFile string

// gpx2JsonCmd represents the gpx2Json command
var gpx2JsonCmd = &cobra.Command{
	Use:   "gpx2Json",
	Short: "gpx file to json",
	Run: func(cmd *cobra.Command, args []string) {
		gpx2Json := pkg.NewGpx2Json(gpx2JsonFile)
		err := gpx2Json.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gpx2JsonCmd)
	gpx2JsonCmd.Flags().StringVarP(&gpx2JsonFile, "file", "f", "", "gpx file import")
}
