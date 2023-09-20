/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"btl/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var gistUrl string

// gistCmd represents the gist command
var gistCmd = &cobra.Command{
	Use:   "gist",
	Short: "download a gist to local md file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"gist_url": gistUrl,
		}).Info("start download gist")

		if gistUrl == "" {
			log.Fatalf("gist url is empty, exit")
		}

		github := pkg.NewGithub()
		filepath, err := github.DownloadGist(gistUrl)
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

	gistCmd.Flags().StringVarP(&gistUrl, "gist_id", "u", "", "gist url")
	_ = gistCmd.MarkFlagRequired("gist_id")
}
