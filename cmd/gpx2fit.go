package cmd

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
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
		log.Infof("gpxFiles size %d", len(gpxFiles))
		var err error
		_ = os.Mkdir("gpx2fit", 0755)

		for index, item := range gpxFiles {
			item := item
			fname := fmt.Sprintf("%d.fit", index+1)

			fitFile := path.Join("gpx2fit", fname)
			err = gpx2fit(item, fitFile)
			if err != nil {
				log.Errorf("gpx2fit failed. file is %s", item)
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
	log.Infof("cmdStr is <%s>", cmdStr)
	out, err := exec.Command("/bin/bash", "-c", cmdStr).CombinedOutput()
	if err != nil {
		return err
	}
	log.Info(string(out))
	return nil
}
