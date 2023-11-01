/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

var gpxFile string

// gpxVideoCmd represents the gpxVideo command
var gpxVideoCmd = &cobra.Command{
	Use:   "gpxVideo",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.NewGpxVideo([]string{gpxFile})
		if err != nil {
			log.Errorf("gpx video error: %+v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gpxVideoCmd)
	gpxVideoCmd.Flags().StringVarP(&gpxFile, "file", "f", "", "xx.gpx")
}
