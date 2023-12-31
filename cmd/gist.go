/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

// gistCmd represents the gist command
var gistCmd = &cobra.Command{
	Use:   "gist",
	Short: "download a gist to local md file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error("need gist_id param. eg: btl gist xxxxx")
			return
		}
		gistId := args[0]
		github := pkg.NewGithub()
		filepath, err := github.DownloadGist(gistId)
		if err != nil {
			log.Fatalf("download gist error %+v", err)
		}
		md := pkg.NewMarkdown(filepath)
		err = md.ParseImages()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gistCmd)
}
