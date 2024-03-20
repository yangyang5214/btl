package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

// xhsCmd represents the xhs command
var xhsCmd = &cobra.Command{
	Use:   "xhs",
	Short: "xhs images",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		xhs := pkg.NewXhs(urlStr)
		err := xhs.Run()
		if err != nil {
			panic(err)
		}
	},
}
var urlStr string

func init() {
	rootCmd.AddCommand(xhsCmd)
	xhsCmd.Flags().StringVarP(&urlStr, "url", "u", "", "https://www.xiaohongshu.com/explore/65cf5f0900000000070278e1")
}
