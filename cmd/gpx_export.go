/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
	"github.com/yangyang5214/btl/pkg/gpx_export"
)

var (
	app      string
	username string
	password string
	outDir   string
)

// gpxExportCmd represents the gpxExport command
var gpxExportCmd = &cobra.Command{
	Use:   "gpx_export",
	Short: "export gpx from apps",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ge := gpx_export.NewGpxExport(app, username, password)

		absPath, err := filepath.Abs(outDir)
		if err != nil {
			panic(err)
		}
		if !pkg.FileExists(absPath) {
			err = os.MkdirAll(absPath, 0755)
			if err != nil {
				panic(err)
			}
		}
		ge.SetExportDir(absPath)
		err = ge.Run()
		if err != nil {
			log.Errorf("export gpx failed: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(gpxExportCmd)
	gpxExportCmd.Flags().StringVarP(&app, "app", "a", "", strings.Join(gpx_export.Apps, "\n"))
	gpxExportCmd.Flags().StringVarP(&username, "user", "u", "", "username")
	gpxExportCmd.Flags().StringVarP(&password, "pwd", "p", "", "password")
	gpxExportCmd.Flags().StringVarP(&outDir, "out", "o", "garmin_export_out", "garmin_export_out dir")

	_ = gpxExportCmd.MarkFlagRequired("app")
	_ = gpxExportCmd.MarkFlagRequired("user")
	_ = gpxExportCmd.MarkFlagRequired("pwd")
}
