/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "btl",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	username string
	password string
	filePath string
)

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&username, "user", "u", "", "username. email")
	rootCmd.PersistentFlags().StringVarP(&password, "pwd", "p", "", "password")
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "input file")
}
