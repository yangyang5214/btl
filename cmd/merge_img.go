package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/merge_img"
)

// mergeImgCmd represents the mergeImg command
var mergeImgCmd = &cobra.Command{
	Use:   "merge_img",
	Short: "merge gpx route images",
	Run: func(cmd *cobra.Command, args []string) {
		mergeImg := merge_img.NewMergeImg(fgFile, bgFile, log.DefaultLogger)
		err := mergeImg.Run("result.png")
		if err != nil {
			panic(err)
		}
	},
}

var (
	bgFile string
	fgFile string
)

func init() {
	rootCmd.AddCommand(mergeImgCmd)
	mergeImgCmd.PersistentFlags().StringVarP(&bgFile, "bg", "", "", "bg image")
	mergeImgCmd.PersistentFlags().StringVarP(&fgFile, "fg", "", "", "fg image")
}
