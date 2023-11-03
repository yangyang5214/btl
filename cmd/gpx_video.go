/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
	"github.com/yangyang5214/btl/pkg/utils"
)

var gpxFile string

// gpxVideoCmd represents the gpxVideo command
var gpxVideoCmd = &cobra.Command{
	Use:   "gpxv",
	Short: "gpx to video",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if gpxFile == "" {
			files = utils.FindGpxFiles(dirPath)
		} else {
			files = []string{gpxFile}
		}
		if len(files) == 0 {
			log.Infof("no gpx file set")
			return
		}

		log.Infof("satrt gpxv, gpx file size is %d", len(files))

		cs, err := parserColor()
		if err != nil {
			log.Info("parse color <%s> error: %v", err)
			return
		}

		gpxv, err := pkg.NewGpxVideo(files, cs)
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
	gpxVideoCmd.Flags().StringVarP(&colorStr, "color", "c", "green", "red green or random")
}
