/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/cmd"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	cmd.Execute()
}
