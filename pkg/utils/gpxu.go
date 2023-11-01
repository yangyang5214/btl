package utils

import (
	"fmt"
	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
	"path/filepath"
	"strconv"
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
