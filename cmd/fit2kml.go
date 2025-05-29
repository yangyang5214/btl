package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

var fit2kmlCmd = &cobra.Command{
	Use: "fit2kml",
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.NewFit2Kml(inputPath, outputPath).Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fit2kmlCmd)
}
