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
	Use:   "gpxv",
	Short: "gpx to video",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if gpxFile == "" {
			log.Infof("no gpx file set")
			return
		}
		gpxv, err := pkg.NewGpxVideo([]string{gpxFile})
		if err != nil {
			log.Errorf("init error: %+v", err)
		}
		err = gpxv.Run()
		if err != nil {
			log.Errorf("gpx video error: %+v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gpxVideoCmd)
	gpxVideoCmd.Flags().StringVarP(&gpxFile, "file", "f", "", "xx.gpx")
}
