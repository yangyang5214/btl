package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/utils"
	"os"
	"os/exec"
	"path"
	"time"
)

var gpx2fitCmd = &cobra.Command{
	Use:   "gpx2fit",
	Short: "gpx 2 fit",
	Run: func(cmd *cobra.Command, args []string) {

		gpxFiles := utils.FindGpxFiles(".")

		var err error
		_ = os.Mkdir("gpx2fit", 0755)
		for _, item := range gpxFiles {
			fitFile := path.Join("gpx2fit", fmt.Sprintf("%d.fit", time.Now().Unix()))
			err = gpx2fit(item, fitFile)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(gpx2fitCmd)
}

func gpx2fit(gpxFile string, fitFile string) error {
	input := path.Join("/tmp", fmt.Sprintf("%d.gpx", time.Now().Unix()))
	defer os.Remove(input)
	err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cp '%s' %s", gpxFile, input)).Run()
	if err != nil {
		return err
	}
	cmdStr := fmt.Sprintf("java -jar ~/gpx2fit.jar %s %s", input, fitFile)
	fmt.Println(cmdStr)
	out, err := exec.Command("/bin/bash", "-c", cmdStr).CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
