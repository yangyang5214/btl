/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg"

	"github.com/spf13/cobra"
)

// csv2mindCmd represents the csv2mind command
var csv2mindCmd = &cobra.Command{
	Use:   "csv2mind",
	Short: "csv file covert to mind map",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("input file is %s", csvFile)
		j := pkg.NewJsonGroup(csvFile, fields)
		err := j.Run()
		if err != nil {
			log.Fatalf("errors in csv2mind: %v", err)
		}
		log.Info("csv2mind successfully converted")
	},
}

var csvFile string

func init() {
	rootCmd.AddCommand(csv2mindCmd)

	csv2mindCmd.Flags().StringVarP(&csvFile, "csv_file", "f", "", "csv file path")
	csv2mindCmd.Flags().StringSliceVarP(&fields, "fields", "g", []string{}, "Fields that require group by")
	_ = csv2mindCmd.MarkFlagRequired("json_file")
	_ = csv2mindCmd.MarkFlagRequired("fields")
}
