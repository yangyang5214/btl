/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

// screenshotCmd represents the screenshot command
var screenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "screenshot for input html",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ss := pkg.NewScreenshot("result.png", htmlPath)
		err := ss.Run()
		if err != nil {
			panic(err)
		}
	},
}

var (
	htmlPath string
)

func init() {
	rootCmd.AddCommand(screenshotCmd)
	screenshotCmd.Flags().StringVar(&htmlPath, "html", "", "index.html")
}
