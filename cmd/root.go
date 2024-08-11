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
	username   string
	password   string
	inputPath  string
	outputPath string
	debug      bool
)

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&username, "user", "u", "", "username. email")
	rootCmd.PersistentFlags().StringVarP(&password, "pwd", "p", "", "password")
	rootCmd.PersistentFlags().StringVarP(&inputPath, "input", "f", "", "input file")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "default not debug")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "", "output file")
}
