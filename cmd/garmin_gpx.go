/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg"

	"github.com/spf13/cobra"
)

// garminGpxCmd represents the garminGpx command
var garminGpxCmd = &cobra.Command{
	Use:   "garminGpx",
	Short: "merge garmin gpx files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("garminGpx cmd called")

		pwd, err := pkg.GetPwd()
		if err != nil {
			log.Fatalf("error getting pwd: %v", err)
		}
		gpx := pkg.GarminGpx{
			CurrentDir: pwd,
		}
		err = gpx.Run()
		if err != nil {
			log.Fatalf("error running garmin gpx: %v", err)
		}

		log.Info("garminGpx merge success")
	},
}

func init() {
	rootCmd.AddCommand(garminGpxCmd)
}
