package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/gpx2video"
)

// gpx2videoCmd represents the gpx2video command
var gpx2videoCmd = &cobra.Command{
	Use:   "gpx2video",
	Short: "gen video by gpx",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "route video",
	Run: func(cmd *cobra.Command, args []string) {
		s := gpx2video.NewGpxVideo(filePath, log.DefaultLogger)
		err := s.Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gpx2videoCmd)
	gpx2videoCmd.AddCommand(routeCmd)
}
