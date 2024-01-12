/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

var (
	statFile string
)

// gpxStatCmd represents the gpxStat command
var gpxStatCmd = &cobra.Command{
	Use:   "gpx_stat",
	Short: "gpx stat",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gs := pkg.NewGpxStat(statFile)
		r, err := gs.Run()
		if err != nil {
			panic(err)
		}
		for _, data := range r {
			log.Info(data.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(gpxStatCmd)
	gpxStatCmd.Flags().StringVarP(&statFile, "file", "f", "", "gpx file")
}
