/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg/gpx_amap"
	"github.com/yangyang5214/btl/pkg/utils"
	"golang.org/x/image/colornames"
	"image/color"

	"github.com/spf13/cobra"
)

// gpxAmapCmd represents the gpxAmap command
var gpxAmapCmd = &cobra.Command{
	Use:   "gpx_amap",
	Short: `gen image by amap`,
	Run: func(cmd *cobra.Command, args []string) {
		files := utils.FindGpxFiles(".")
		gamap := gpx_amap.NewGpxAmap(files, "result.png")
		gamap.SetColors([]color.Color{
			colornames.Red,
		})
		err := gamap.Run()
		if err != nil {
			log.Errorf("run gamap failed: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(gpxAmapCmd)
}
