package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/yangyang5214/btl/pkg/amap_replaying"
)

// amapReplayingCmd represents the amapReplaying command
var amapReplayingCmd = &cobra.Command{
	Use:   "amap_replaying",
	Short: "amap replay",
	Long:  `https://lbs.amap.com/demo/javascript-api-v2/example/marker/replaying-historical-running-data`,
	Run: func(cmd *cobra.Command, args []string) {
		gpxData, err := gpx.ParseFile(gpxFile)
		if err != nil {
			panic(err)
		}
		ar := amap_replaying.NewAmapReplaying(gpxData, log.DefaultLogger)
		err = ar.Run()
		if err != nil {
			panic(err)
		}
	},
}

var gpxFile string

func init() {
	rootCmd.AddCommand(amapReplayingCmd)
	amapReplayingCmd.Flags().StringVarP(&gpxFile, "gpx_file", "f", "", "a gpx file")
}
