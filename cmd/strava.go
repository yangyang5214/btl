package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/strava"
)

// stravaCmd represents the strava command
var stravaCmd = &cobra.Command{
	Use:   "strava",
	Short: "strava upload file",
	Run: func(cmd *cobra.Command, args []string) {
		err := strava.NewUploadFit(filePath, log.DefaultLogger).Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(stravaCmd)
}
