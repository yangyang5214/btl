package pkg

import (
	"fmt"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
	"path/filepath"
	"strconv"
)

type GpxMap struct {
	files       []string
	smCtx       *sm.Context
	attribution string
}

func NewGpxMap(files []string, attribution string) *GpxMap {
	return &GpxMap{
		files:       files,
		smCtx:       sm.NewContext(),
		attribution: attribution,
	}
}

func (g *GpxMap) getWidthHeight(positions []s2.LatLng) (int, int) {
	southPoint, northPoint := positions[0], positions[0]
	for _, position := range positions {
		if position.Lat < southPoint.Lat {
			southPoint = position
		}
		if position.Lat > northPoint.Lat {
			northPoint = position
		}
	}

	log.Infof("northPoint: %s", fmt.Sprintf("%s,%s", northPoint.Lng, northPoint.Lat))
	log.Infof("southPoint: %s", fmt.Sprintf("%s,%s", southPoint.Lng, southPoint.Lat))

	south, _ := strconv.ParseFloat(southPoint.Lat.String(), 10)
	north, _ := strconv.ParseFloat(northPoint.Lat.String(), 10)
	height := (north - south) * 1000 / 4
	if height < 600 {
		return 800, 600
	}
	return int(height * 1.5), int(height)
}

func (g *GpxMap) parsePositions() ([]s2.LatLng, error) {
	var positions []s2.LatLng
	for _, f := range g.files {
		p, err := filepath.Abs(f)
		if err != nil {
			return nil, err
		}
		gpxData, err := gpx.ParseFile(p)
		if err != nil {
			return nil, err
		}
		for _, trk := range gpxData.Tracks {
			for _, seg := range trk.Segments {
				for _, pt := range seg.Points {
					positions = append(positions, s2.LatLngFromDegrees(pt.GetLatitude(), pt.GetLongitude()))
				}
			}
		}
	}
	return positions, nil
}

func (g *GpxMap) Run(imgPath string) error {
	color, _ := sm.ParseColorString("red")
	positions, err := g.parsePositions()
	if err != nil {
		return err
	}
	width, height := g.getWidthHeight(positions)
	log.Infof("use height=%d, width=%d", height, width)
	g.smCtx.SetSize(width, height)
	g.smCtx.AddObject(sm.NewPath(positions, color, 3))

	titleProvider := sm.NewTileProviderOpenStreetMaps()
	titleProvider.Attribution = g.attribution
	g.smCtx.SetTileProvider(titleProvider)
	img, err := g.smCtx.Render()
	if err != nil {
		return err
	}

	if err = gg.SavePNG(imgPath, img); err != nil {
		return err
	}
	return nil
}
