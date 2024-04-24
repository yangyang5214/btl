package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/yangyang5214/btl/pkg"

	"github.com/spf13/cobra"
)

// gpxSpeedCmd represents the gpxSpeed command
var gpxSpeedCmd = &cobra.Command{
	Use:   "gpx_speed",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.NewGpxSpeed(filePath, speed, log.DefaultLogger).Run()
		if err != nil {
			panic(err)
		}
	},
}

var (
	speed int32
)

func init() {
	rootCmd.AddCommand(gpxSpeedCmd)
	gpxSpeedCmd.Flags().Int32Var(&speed, "speed", 100, "default set 100 km/h")
}
