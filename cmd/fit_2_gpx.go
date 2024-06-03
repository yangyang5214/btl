package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

// fitToGpxCmd represents the fitToGpx command
var fitToGpxCmd = &cobra.Command{
	Use:   "fit2gpx",
	Short: "fit to gpx file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fit := pkg.NewFit2Gpx(filePath, log.DefaultLogger)
		err := fit.Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fitToGpxCmd)
}
