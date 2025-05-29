/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/yangyang5214/btl/pkg"

	"github.com/spf13/cobra"
)

// fitTrimCmd represents the fitTrim command
var fitTrimCmd = &cobra.Command{
	Use: "fit_trim",
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.NewFitTrim(inputPath)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fitTrimCmd)
}
