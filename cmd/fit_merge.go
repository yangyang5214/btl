package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

var fitMergeCmd = &cobra.Command{
	Use:   "fit_merge",
	Short: "merge current dir fits",
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.NewFitMerge().Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fitMergeCmd)
}
