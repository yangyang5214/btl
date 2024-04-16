package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/yangyang5214/btl/pkg/gpx_export"

	"github.com/spf13/cobra"
)

// garminUploadCmd represents the garminUpload command
var garminUploadCmd = &cobra.Command{
	Use:   "garmin_upload",
	Short: "garmin_upload",
	Long:  `upload gpx/fit file to Garmin.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := gpx_export.NewGarminUpload(username, password, filePath, isCn, log.DefaultLogger).Run()
		if err != nil {
			panic(err)
		}
	},
}

var (
	isCn bool
)

func init() {
	rootCmd.AddCommand(garminUploadCmd)
	garminUploadCmd.Flags().BoolVarP(&isCn, "is_cn", "c", true, "is garmin.cn.国内版")
}
