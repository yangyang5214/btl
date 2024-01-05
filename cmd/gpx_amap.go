/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/gpx_amap"
	"github.com/yangyang5214/btl/pkg/utils"
)

var (
	files        []string
	amapStyle    string
	ValidChoices = []string{"whitesmoke", "grey", "dark", "light", "fresh", "blue", "darkblue", "macaron"}
)

// gpxAmapCmd represents the gpxAmap command
var gpxAmapCmd = &cobra.Command{
	Use:   "gpx_amap",
	Short: `gen image by amap`,
	Run: func(cmd *cobra.Command, args []string) {
		validateMyFlag()
		if len(files) == 0 {
			files = utils.FindGpxFiles(".")
		}
		gamap := gpx_amap.NewGpxAmap(amapStyle)
		gamap.SetFiles(files)
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
	gpxAmapCmd.Flags().StringVarP(&amapStyle, "style", "s", "", strings.Join(ValidChoices, "\n"))
}

func validateMyFlag() {
	for _, choice := range ValidChoices {
		if amapStyle == choice {
			return
		}
	}
	log.Infof("Invalid value for --myflag. Allowed values are: %v\n", ValidChoices)
	os.Exit(1)
}
