package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

// gpx2kmlCmd represents the gpx2kml command
var gpx2kmlCmd = &cobra.Command{
	Use:   "gpx2kml",
	Short: "gpx to kml file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		opts := pkg.WithSpeed(speed)
		err := pkg.NewGpx2Kml(filePath, log.DefaultLogger, opts).Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gpx2kmlCmd)
	gpx2kmlCmd.Flags().Int32Var(&speed, "speed", 100, "set speed")
}
