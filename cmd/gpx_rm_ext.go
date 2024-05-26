/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tkrajina/gpxgo/gpx"
	"os"
)

// gpxRmExt represents the gpxRemoveExtension command
var gpxRmExt = &cobra.Command{
	Use:   "gpx_rm_extension",
	Short: "gpx remove extension",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := Main()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gpxRmExt)
}

func Main() error {
	newGpxData, err := gpx.ParseFile(filePath)
	if err != nil {
		return err
	}

	var tracks []gpx.GPXTrack
	for _, track := range newGpxData.Tracks {
		var segments []gpx.GPXTrackSegment
		for _, segment := range track.Segments {
			var newPoints []gpx.GPXPoint
			for _, point := range segment.Points {
				point.Extensions = gpx.Extension{}
				newPoints = append(newPoints, point)
			}
			segment.Points = newPoints
			segments = append(segments, segment)
		}
		track.Segments = segments
		tracks = append(tracks, track)
	}

	newGpxData.Tracks = tracks
	newXml, err := newGpxData.ToXml(gpx.ToXmlParams{
		Indent: true,
	})
	f1, err := os.Create("remove_extension_result.gpx")
	if err != nil {
		return err
	}
	defer f1.Close()
	_, _ = f1.Write(newXml)
	return nil
}
