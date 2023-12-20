/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/gpx_amap"
	"github.com/yangyang5214/btl/pkg/utils"
)

var (
	files []string
)

// gpxAmapCmd represents the gpxAmap command
var gpxAmapCmd = &cobra.Command{
	Use:   "gpx_amap",
	Short: `gen image by amap`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(files) == 0 {
			files = utils.FindGpxFiles(".")
		}
		gamap := gpx_amap.NewGpxAmap(files)
		err := gamap.Run()
		if err != nil {
			log.Errorf("run gamap failed: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(gpxAmapCmd)
	gpxAmapCmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "xx.gpx")
}
