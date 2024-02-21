/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log2 "github.com/go-kratos/kratos/v2/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/gpx_amap"
	"github.com/yangyang5214/btl/pkg/utils"
)

var (
	files     []string
	amapStyle string
)

// gpxAmapCmd represents the gpxAmap command
var gpxAmapCmd = &cobra.Command{
	Use:   "gpx_amap",
	Short: `gen image by amap`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(files) == 0 {
			files = utils.FindGpxFiles(".")
		}
		gamap := gpx_amap.NewGpxAmap(amapStyle, log2.DefaultLogger)
		gamap.SetFiles(files)
		gamap.Screenshot()
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
	gpxAmapCmd.Flags().StringVarP(&amapStyle, "style", "s", "8ee61a45840f14ac60f33a799fbd00d8", "amap style id")
}
