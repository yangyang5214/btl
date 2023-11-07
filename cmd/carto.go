/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strconv"
	"strings"

	"github.com/golang/geo/s2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

var (
	cartoName string
	location  string
	km        int
)

// cartoCmd represents the carto command
var cartoCmd = &cobra.Command{
	Use:   "carto",
	Short: "download carto map",
	Long:  `default zoom range [8,15]`,
	Run: func(cmd *cobra.Command, args []string) {
		point := strings.Split(location, ",")
		if len(point) != 2 {
			log.Infof("Invalid carto location")
			return
		}
		lat, err := strconv.ParseFloat(point[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		lng, err := strconv.ParseFloat(point[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		center := s2.LatLngFromDegrees(lat, lng)
		for i := 8; i <= 13; i++ {
			carto := pkg.NewCarto(cartoName, i, center, float64(km))
			err = carto.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cartoCmd)

	//https://github.com/CartoDB/basemap-styles
	cartoCmd.Flags().StringVarP(&cartoName, "carto_name", "n", "carto-dark", gpxMapUsage())
	cartoCmd.Flags().StringVarP(&location, "location", "l", "121.455708,31.249574", "lat,lng. \n  https://lbs.amap.com/tools/picker)")
	cartoCmd.Flags().IntVarP(&km, "km", "k", 300, "")
}
