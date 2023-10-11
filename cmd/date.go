/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// dateCmd represents the date command
var dateCmd = &cobra.Command{
	Use:   "date",
	Short: "date format",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			defaultFunc()
			return
		}
		param := args[0]
		if len(param) == 10 || len(param) == 13 {
			formatTimestamp(param)
		} else {
			defaultFunc()
		}
	},
}

func defaultFunc() {
	now := time.Now()
	log.Infof("Now date: %s", now.Format(time.DateTime))
	log.Infof("Now timestamp: %d", now.UnixMilli())
}

func formatTimestamp(timestamp string) {
	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Error(err)
		return
	}
	if len(timestamp) == 13 {
		t = t / 1000
	}
	dateStr := time.Unix(t, 0).Format(time.DateTime)
	log.Infof(dateStr)
}

func init() {
	rootCmd.AddCommand(dateCmd)
}
