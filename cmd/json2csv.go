/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

// json2csvCmd represents the json2csv command
var json2csvCmd = &cobra.Command{
	Use:   "json2csv",
	Short: "json to csv",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		j2c := pkg.NewJson2Csv(fPath)
		err := j2c.Run()
		if err != nil {
			log.Errorf("json to csv error: %+v", err)
		}
	},
}

var fPath string

func init() {
	rootCmd.AddCommand(json2csvCmd)
}
