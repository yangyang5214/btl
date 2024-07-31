package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/utils"
	"os"
	"os/exec"
)

var gpx2fitCmd = &cobra.Command{
	Use:   "gpx2fit",
	Short: "gpx 2 fit",
	Run: func(cmd *cobra.Command, args []string) {

		gpxFiles := utils.FindGpxFiles(".")

		err := os.Mkdir("gpx2fit", 0755)
		if err != nil {
			panic(err)
		}
		for index, gpxFile := range gpxFiles {
			resultFile := fmt.Sprintf("%d.fit", index)
			cmdStr := fmt.Sprintf("java -jar ~/gpx2fit.jar %s %s", gpxFile, resultFile)
			out, err := exec.Command("/bin/bash", "-c", cmdStr).CombinedOutput()
			if err != nil {
				panic(err)
			}
			fmt.Println(string(out))
		}
	},
}

func init() {
	rootCmd.AddCommand(gpx2fitCmd)
}
