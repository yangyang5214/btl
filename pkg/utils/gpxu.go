package utils

import (
	"fmt"
	. "image/color"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
)

var (
	DefaultColors = []Color{
		//RGBA{A: 0xff},                 //black
		RGBA{B: 0xff, A: 0xff},          //blue
		RGBA{G: 0xff, A: 0xff},          //green
		RGBA{R: 0xff, G: 0x7f, A: 0xff}, //orange
		RGBA{R: 0xff, A: 0xff},          //red
		RGBA{R: 0xff, G: 0xff, A: 0xff}, //yellow
	}
)

func ParseGpxData(files []string) ([]*gpx.GPX, error) {
	var results []*gpx.GPX
	for _, f := range files {
		p, err := filepath.Abs(f)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		gpxData, err := gpx.ParseFile(p)
		if err != nil {
			log.Errorf("gpx parse file <%s> error: %v", p, err)
			return nil, errors.WithStack(err)
		}
		results = append(results, gpxData)
	}
	return results, nil
}

func ParsePositions(datas []*gpx.GPX) ([][]s2.LatLng, error) {
	var positions [][]s2.LatLng
	for _, gpxData := range datas {
		var local []s2.LatLng
		for _, trk := range gpxData.Tracks {
			for _, seg := range trk.Segments {
				for _, pt := range seg.Points {
					local = append(local, s2.LatLngFromDegrees(pt.GetLatitude(), pt.GetLongitude()))
				}
			}
		}
		positions = append(positions, local)
	}
	return positions, nil
}

func GenWidthHeight(positions [][]s2.LatLng) (int, int) {
	var maxPointSize int
	southPoint, northPoint := positions[0][0], positions[0][0]
	for _, points := range positions {
		for _, point := range points {
			if point.Lat < southPoint.Lat {
				southPoint = point
			}
			if point.Lat > northPoint.Lat {
				northPoint = point
			}
		}
		if len(points) > maxPointSize {
			maxPointSize = len(points)
		}
	}

	log.Infof("northPoint: %s", fmt.Sprintf("%s,%s", northPoint.Lng, northPoint.Lat))
	log.Infof("southPoint: %s", fmt.Sprintf("%s,%s", southPoint.Lng, southPoint.Lat))

	south, _ := strconv.ParseFloat(southPoint.Lat.String(), 10)
	north, _ := strconv.ParseFloat(northPoint.Lat.String(), 10)
	height := (north - south) * 1000 / 4
	if height < 600 {
		if maxPointSize > 2000 {
			return 1028, 1000
		} else {
			return 800, 600
		}
	}
	return int(height * 1.5), int(height)
}

func CountPoints(positions [][]s2.LatLng) int {
	var r int
	for _, sub := range positions {
		r = r + len(sub)
	}
	return r
}

func FindGpxFiles(dirPath string) []string {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil
	}
	d, err := os.ReadDir(absPath)
	if err != nil {
		return nil
	}
	var r []string
	for _, item := range d {
		name := item.Name()
		p := path.Join(absPath, name)
		if item.IsDir() {
			r = append(r, FindGpxFiles(p)...)
		}
		if strings.HasSuffix(name, ".gpx") {
			r = append(r, p)
		}
	}
	return r
}

func GetColor(index int, colors []Color) Color {
	return colors[index%len(colors)]
}
