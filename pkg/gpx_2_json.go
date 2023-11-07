package pkg

import (
	"encoding/json"
	"os"

	"github.com/tkrajina/gpxgo/gpx"
)

type Gpx2Json struct {
	gpxFile string
}

func NewGpx2Json(gpxFile string) *Gpx2Json {
	return &Gpx2Json{gpxFile: gpxFile}
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (g *Gpx2Json) Run() error {
	gpxData, err := gpx.ParseFile(g.gpxFile)
	if err != nil {
		return err
	}

	var points []*LatLng
	for _, track := range gpxData.Tracks {
		for _, seg := range track.Segments {
			for _, pt := range seg.Points {
				points = append(points, &LatLng{
					Lat: pt.GetLatitude(),
					Lng: pt.GetLongitude(),
				})
			}
		}
	}
	result, err := os.Create("result.json")
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(points)
	if err != nil {
		return err
	}
	_, err = result.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
