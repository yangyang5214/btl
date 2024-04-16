/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

// json2mindCmd represents the json2mind command
var json2mindCmd = &cobra.Command{
	Use:   "json2mind",
	Short: "json lines file transform to md file, then convert to mind map",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		j := pkg.NewJsonGroup(filePath, fields)
		err := j.Run()
		if err != nil {
			log.Fatalf("errors in json2mind: %v", err)
		}
		log.Info("json2mind successfully converted")
	},
}

var fields []string

func init() {
	rootCmd.AddCommand(json2mindCmd)
	json2mindCmd.Flags().StringSliceVarP(&fields, "fields", "g", []string{}, "Fields that require group by")
	_ = json2mindCmd.MarkFlagRequired("fields")
}
