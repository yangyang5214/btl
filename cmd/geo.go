/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
	"os"
	"strings"
	"time"
)

// geoCmd represents the geo command
var geoCmd = &cobra.Command{
	// https://lbs.amap.com/api/webservice/guide/api/georegeo
	Use:   "geo",
	Short: "高德 geo parser",
	Run: func(cmd *cobra.Command, args []string) {
		logHelper := log.NewHelper(log.DefaultLogger)
		if address == "" && inputPath == "" {
			logHelper.Error("need address param. eg: btl geo -a '上海市浦东新区陆家嘴街道'")
			return
		}

		amap := pkg.NewAmap(logHelper)
		if address != "" {
			location, err := amap.GetLocationByAddress(address)
			if err != nil {
				panic(err)
			}
			logHelper.Infof("location: %v", location)
		} else {
			err := processCSV(logHelper, amap)
			if err != nil {
				panic(err)
			}
		}
	},
}

func processCSV(logHelper *log.Helper, amap *pkg.Amap) error {
	fileHandle, err := os.Open(inputPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer fileHandle.Close()

	reader := csv.NewReader(fileHandle)
	records, err := reader.ReadAll()
	if err != nil {
		return errors.WithStack(err)
	}

	resultFile, err := os.Create(fmt.Sprintf("geo_result_%d.csv", time.Now().UnixMilli()))
	if err != nil {
		return errors.WithStack(err)
	}
	defer resultFile.Close()

	writer := csv.NewWriter(resultFile)
	defer writer.Flush()

	var location string

	for _, record := range records {
		//set index
		addressLine := record[csvIndex]

		logHelper.Infof("address is: %v", addressLine)

		// Simulate API call to get location based on address
		time.Sleep(1 * time.Second) // Rate limit

		var long, lat string

		location, err = amap.GetLocationByAddress(addressLine)
		if err != nil {
			logHelper.Error(err) //ignore
		} else {
			longLat := strings.Split(location, ",")
			long = longLat[0]
			lat = longLat[1]
		}

		record = append(record, long, lat)

		// Write the updated record to the result file
		err = writer.Write(record)
		if err != nil {
			logHelper.Errorf("Failed to write record %v: %v", record, err)
			return err
		}
	}
	return nil
}

var (
	address  string
	csvIndex int
)

func init() {
	rootCmd.AddCommand(geoCmd)
	geoCmd.Flags().StringVarP(&address, "address", "a", "", "set address")
	geoCmd.Flags().IntVarP(&csvIndex, "csv_index", "i", 0, "csv index for address")
}
