/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg"

	"github.com/spf13/cobra"
)

// gpxSplitCmd represents the gpxSplit command
var gpxSplitCmd = &cobra.Command{
	Use:   "gpx_zip",
	Short: "zip gpx file to limited size.(remove some points)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if s == 0 {
			s = 25
		}
		if f == "" {
			log.Info("filepath is required")
			return
		}
		gs := pkg.NewGpxZip(f, s)
		err := gs.Run()
		if err != nil {
			panic(err)
		}
	},
}

var (
	f string
	s int32
)

func init() {
	rootCmd.AddCommand(gpxSplitCmd)
	gpxSplitCmd.Flags().Int32("size", s, "limit size(MB)")
}
