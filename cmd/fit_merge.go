package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
	"os"
	"path"
)

var fitMergeCmd = &cobra.Command{
	Use:   "fit_merge",
	Short: "merge current dir fits",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		p := path.Join(homeDir, "merge-fit.jar")
		err = pkg.NewFitMerge("java -jar " + p).Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fitMergeCmd)
}
