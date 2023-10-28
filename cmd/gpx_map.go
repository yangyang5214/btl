package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
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
		gpxMap := pkg.NewGpxMap(files, attribution, titleName)
		resultImg := "result.png"
		err := gpxMap.Run(resultImg)
		if err != nil {
			log.Errorf("run  gpxMap error: %+v", err)
			return
		}
		log.Infof("result image is: %s", resultImg)
	},
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
		if item.IsDir() {
			continue
		}
		name := item.Name()
		if strings.HasSuffix(name, ".gpx") {
			r = append(r, path.Join(absPath, name))
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
}
