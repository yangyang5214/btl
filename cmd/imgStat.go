package cmd

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

// imgStatCmd represents the imgStat command
var imgStatCmd = &cobra.Command{
	Use:   "img_stat",
	Short: "Gen stat in image",
	Run: func(cmd *cobra.Command, args []string) {
		imgStat := pkg.NewImgStat(filePath, log.DefaultLogger)

		err := imgStat.Run(initInfos())
		if err != nil {
			panic(err)
		}
	},
}

func initInfos() (infos []*pkg.StatInfo) {
	infos = append(infos, &pkg.StatInfo{
		Label: "距离",
		Value: "10 KM",
	})
	infos = append(infos, &pkg.StatInfo{
		Label: "时间",
		Value: "1h 10m",
	})
	infos = append(infos, &pkg.StatInfo{
		Label: "最大速度",
		Value: "30 km/h",
	})
	return
}

func init() {
	rootCmd.AddCommand(imgStatCmd)
}
