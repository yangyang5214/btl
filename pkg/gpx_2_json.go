package pkg

import (
	"encoding/json"
	"os"

	"github.com/yangyang5214/btl/pkg/utils"
)

type Gpx2Json struct {
	gpxFile string
}

func NewGpx2Json(gpxFile string) *Gpx2Json {
	return &Gpx2Json{gpxFile: gpxFile}
}

func (g *Gpx2Json) Run() error {
	points, err := utils.GetPoints([]string{g.gpxFile})
	if err != nil {
		return err
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
