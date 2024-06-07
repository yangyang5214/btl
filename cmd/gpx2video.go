package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/gpx2video"
	"os"
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
		tmpDir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = os.Remove(tmpDir)
		}()
		s, err := gpx2video.NewRouteVideo(filePath, log.DefaultLogger, tmpDir)
		if err != nil {
			panic(err)
		}
		err = s.Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gpx2videoCmd)
	gpx2videoCmd.AddCommand(routeCmd)
}
