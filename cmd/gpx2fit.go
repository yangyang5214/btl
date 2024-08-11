package cmd

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/yangyang5214/btl/pkg"
	"github.com/yangyang5214/btl/pkg/utils"
	"os"
	"path"
)

var gpx2fitCmd = &cobra.Command{
	Use:   "gpx2fit",
	Short: "gpx 2 fit",
	Run: func(cmd *cobra.Command, args []string) {
		if inputPath != "" {
			//single
			out := "result.fit"
			if outputPath != "" {
				out = outputPath
			}
			err := gpx2fit(inputPath, out, debug)
			if err != nil {
				panic(err)
			}
			return
		}
		gpxFiles := utils.FindGpxFiles(".")
		log.Infof("gpxFiles size %d", len(gpxFiles))
		var err error
		_ = os.Mkdir("gpx2fit", 0755)

		for index, item := range gpxFiles {
			item := item
			fname := fmt.Sprintf("%d.fit", index+1)
			fitFile := path.Join("gpx2fit", fname)
			err = gpx2fit(item, fitFile, debug)
			if err != nil {
				log.Errorf("gpx2fit failed. file is %s", item)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(gpx2fitCmd)
}

func gpx2fit(gpxFile string, fitFile string, debug bool) error {
	gpxData, err := gpx.ParseFile(gpxFile)
	if err != nil {
		return err
	}
	newXml, err := gpxData.ToXml(gpx.ToXmlParams{
		Indent: true,
	})
	if err != nil {
		return err
	}
	return pkg.GenFitFile(debug, "java -jar ~/gpx2fit.jar", log.NewHelper(log.DefaultLogger), newXml, fitFile)
}
