package cmd

import (
	. "image/color"
	"strings"

	"github.com/yangyang5214/btl/pkg/utils"

	sm "github.com/flopp/go-staticmaps"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

var files []string
var (
	dirPath     string
	attribution string
	titleName   string
	colorStr    string
)

// gpxMapCmd represents the gpxMap command
var gpxMapCmd = &cobra.Command{
	Use:   "gpxm",
	Short: "show gpx in map",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(files) == 0 {
			files = utils.FindGpxFiles(dirPath)
		}
		if len(files) == 0 {
			log.Info("inout gpx files is empty")
			return
		}
		c, err := parserColor()
		if err != nil {
			log.Info("parse color <%s> error: %v", err)
			return
		}
		gpxMap := pkg.NewGpxMap(files, attribution, titleName, c)
		resultImg := "result.png"
		err = gpxMap.Run(resultImg)
		if err != nil {
			log.Errorf("run  gpxMap error: %+v", err)
			return
		}
		log.Infof("result image is: %s", resultImg)
	},
}

func parserColor() ([]Color, error) {
	if colorStr == "random" {
		return []Color{
			//RGBA{A: 0xff},                  //black
			RGBA{B: 0xff, A: 0xff},          //blue
			RGBA{G: 0xff, A: 0xff},          //green
			RGBA{R: 0xff, G: 0x7f, A: 0xff}, //orange
			RGBA{R: 0xff, A: 0xff},          //red
			RGBA{R: 0xff, G: 0xff, A: 0xff}, //yellow
		}, nil
	}
	c, err := sm.ParseColorString(colorStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return []Color{c}, nil
}

func gpxMapUsage() string {
	r := []string{
		"carto-dark(暗黑)",
		"carto-light(白色)",
		"wikimedia(纯线路)",
		"osm",
	}
	return strings.Join(r, "\n")
}

func init() {
	rootCmd.AddCommand(gpxMapCmd)
	gpxMapCmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "xx.gpx")
	gpxMapCmd.Flags().StringVarP(&dirPath, "dir", "d", ".", "")
	gpxMapCmd.Flags().StringVarP(&attribution, "attribution", "a", "", "")
	gpxMapCmd.Flags().StringVarP(&titleName, "name", "n", "carto-dark", gpxMapUsage())
	gpxMapCmd.Flags().StringVarP(&colorStr, "color", "c", "red", "red or green or random")
}
