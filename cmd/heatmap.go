package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg/heatmap"
)

// heatmapCmd represents the heatmap command
var heatmapCmd = &cobra.Command{
	Use:   "heatmap",
	Short: "heatmap for merge_gpx",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		hm := heatmap.NewHeatMap(userId)
		err := hm.Run()
		if err != nil {
			panic(err)
		}
	},
}

var (
	userId string
)

func init() {
	rootCmd.AddCommand(heatmapCmd)
	heatmapCmd.Flags().StringVarP(&userId, "user", "u", "", "merge_gpx user_Id")
}
