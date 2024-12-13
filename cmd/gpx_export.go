package cmd

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
	"github.com/yangyang5214/btl/pkg/gpx_export"
)

var (
	app    string
	outDir string
)

// gpxExportCmd represents the gpxExport command
var gpxExportCmd = &cobra.Command{
	Use:   "gpx_export",
	Short: "export gpx from apps",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ge := gpx_export.NewGpxExport(log.DefaultLogger, app, username, password)

		absPath, err := filepath.Abs(path.Join(outDir))
		if err != nil {
			panic(err)
		}

		out := path.Join(absPath, app)
		if !pkg.FileExists(out) {
			err = os.MkdirAll(out, 0755)
			if err != nil {
				panic(err)
			}
		}
		ge.SetExportDir(out)
		err = ge.Run(true)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gpxExportCmd)
	gpxExportCmd.Flags().StringVarP(&app, "app", "a", "", strings.Join(gpx_export.Apps, "\n"))
	gpxExportCmd.Flags().StringVarP(&outDir, "out", "", "gpx_export_out", "garmin_export_out_dir")

	_ = gpxExportCmd.MarkFlagRequired("app")
	_ = gpxExportCmd.MarkFlagRequired("user")
	_ = gpxExportCmd.MarkFlagRequired("pwd")
}
