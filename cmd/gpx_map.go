package cmd

import (
	sm "github.com/flopp/go-staticmaps"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	Use:   "gpx_map",
	Short: "show gpx in map",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(files) == 0 {
			files = getFiles(dirPath)
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

func parserColor() (color.Color, error) {
	c, err := sm.ParseColorString(colorStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return c, nil
}

func getFiles(dirPath string) []string {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil
	}
	d, err := os.ReadDir(absPath)
	if err != nil {
		return nil
	}
	var r []string
	for _, item := range d {
		name := item.Name()
		p := path.Join(absPath, name)
		if item.IsDir() {
			r = append(r, getFiles(p)...)
		}
		if strings.HasSuffix(name, ".gpx") {
			r = append(r, p)
		}
	}
	return r
}

func gpxMapUsage() string {
	r := []string{
		"carto-dark(暗黑)",
		"carto-light(白色)",
		"wikimedia(纯线路)",
	}
	return strings.Join(r, "\n")
}

func init() {
	rootCmd.AddCommand(gpxMapCmd)
	gpxMapCmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "xx.gpx")
	gpxMapCmd.Flags().StringVarP(&dirPath, "dir", "d", ".", "")
	gpxMapCmd.Flags().StringVarP(&attribution, "attribution", "a", "", "")
	gpxMapCmd.Flags().StringVarP(&titleName, "name", "n", "carto-light", gpxMapUsage())
	gpxMapCmd.Flags().StringVarP(&colorStr, "color", "c", "red", "set color")
}
