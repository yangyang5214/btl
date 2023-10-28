package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

var files []string
var attribution string

// gpxMapCmd represents the gpxMap command
var gpxMapCmd = &cobra.Command{
	Use:   "gpxMap",
	Short: "show gpx in map",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(files) == 0 {
			log.Info("inout gpx files is empty")
			return
		}
		gpxMap := pkg.NewGpxMap(files, attribution)
		resultImg := "result.png"
		err := gpxMap.Run(resultImg)
		if err != nil {
			log.Errorf("run  gpxMap error: %+v", err)
			return
		}
		log.Infof("result image is: %s", resultImg)
	},
}

func init() {
	rootCmd.AddCommand(gpxMapCmd)
	gpxMapCmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "xx.gpx")
	gpxMapCmd.Flags().StringVarP(&attribution, "attribution", "a", "", "beer")
	_ = gpxMapCmd.MarkFlagRequired("files")
}
