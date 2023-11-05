package utils

import (
	sm "github.com/yangyang5214/go-staticmaps"
	"golang.org/x/image/colornames"
	. "image/color"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
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
	smCtx := sm.NewContext()
	for _, position := range positions {
		smCtx.AddObject(sm.NewPath(position, colornames.Yellow, 1))
	}
	bounds := smCtx.DetermineBounds()
	maxLength := bounds.Lat.Length()
	if bounds.Lng.Length() > maxLength {
		maxLength = bounds.Lng.Length()
	}

	width := 0
	dis := int(maxLength * 1_0000)
	gap := dis / 50
	width = (gap)*200 + 800

	log.Printf("dis %v, width %v,", dis, width)

	return width, int(float64(width) * 0.8)
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
