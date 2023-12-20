/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/gpx_export"
)

var (
	app string
)

// gpxExportCmd represents the gpxExport command
var gpxExportCmd = &cobra.Command{
	Use:   "gpx_export",
	Short: "export gpx from apps",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ge := gpx_export.NewGpxExport(app)
		err := ge.Run()
		if err != nil {
			log.Errorf("export gpx failed: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(gpxExportCmd)
	gpxExportCmd.Flags().StringVarP(&app, "app", "a", "", "App name from export")
}
