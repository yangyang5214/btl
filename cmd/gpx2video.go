package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/yangyang5214/btl/pkg/gpx2video"
	"os"
	"strings"
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
		gpxData, err := gpx.ParseFile(filePath)
		if err != nil {
			panic(err)
		}
		s := gpx2video.NewRouteVideo(gpxData, log.DefaultLogger, tmpDir)
		err = s.Run()
		if err != nil {
			panic(err)
		}
	},
}

var imgCmd = &cobra.Command{
	Use:   "img",
	Short: "img overview",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			session *gpx2video.Session
			err     error
			logger  = log.DefaultLogger
		)

		if strings.HasSuffix(filePath, ".gpx") {
			gpxData, err := gpx.ParseFile(filePath)
			if err != nil {
				panic(err)
			}
			session, err = gpx2video.ParseGPX(gpxData)
		} else if strings.HasSuffix(filePath, ".fit") {
			session, err = gpx2video.ParseFit(filePath, log.NewHelper(logger))
		} else {
			panic("不支持的文件类型")
		}

		s := gpx2video.NewImgOverview(session, logger)
		err = s.Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gpx2videoCmd)
	gpx2videoCmd.AddCommand(routeCmd)
	gpx2videoCmd.AddCommand(imgCmd)
}
