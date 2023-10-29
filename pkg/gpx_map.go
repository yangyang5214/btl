package pkg

import (
	"fmt"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
	ic "image/color"
	"path/filepath"
	"strconv"
)

type GpxMap struct {
	files         []string
	smCtx         *sm.Context
	attribution   string
	titleName     string
	tileProviders map[string]*sm.TileProvider
	color         ic.Color
}

func NewGpxMap(files []string, attribution, titleName string, color ic.Color) *GpxMap {
	return &GpxMap{
		files:         files,
		smCtx:         sm.NewContext(),
		attribution:   attribution,
		titleName:     titleName,
		tileProviders: sm.GetTileProviders(),
		color:         color,
	}
}

func (g *GpxMap) getWidthHeight(positions [][]s2.LatLng) (int, int) {
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

func (g *GpxMap) parsePositions() ([][]s2.LatLng, error) {
	var positions [][]s2.LatLng
	for _, f := range g.files {
		p, err := filepath.Abs(f)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		gpxData, err := gpx.ParseFile(p)
		if err != nil {
			log.Errorf("gpx parse file <%s> error: %v", p, err)
			return nil, errors.WithStack(err)
		}
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

func (g *GpxMap) getWeight(post []s2.LatLng) float64 {
	var weight float64
	defer func() {
		log.Infof("points size is %d, line weight %v", len(post), weight)
	}()
	size := len(post)
	weight = float64(size/10_000) + 2.0
	return weight
}

func (g *GpxMap) Run(imgPath string) error {
	positions, err := g.parsePositions()
	if err != nil {
		return errors.WithStack(err)
	}
	width, height := g.getWidthHeight(positions)
	log.Infof("use height=%d, width=%d", height, width)
	g.smCtx.SetSize(width, height)

	for _, post := range positions {
		weight := g.getWeight(post)
		if height <= 1000 {
			weight = 2
		}
		g.smCtx.AddObject(sm.NewPath(post, g.color, weight))
	}

	titleProvider, ok := g.tileProviders[g.titleName]
	if !ok {
		titleProvider = sm.NewTileProviderOpenStreetMaps()
	}

	titleProvider.Attribution = g.attribution

	g.smCtx.SetTileProvider(titleProvider)
	img, err := g.smCtx.Render()
	if err != nil {
		return errors.WithStack(err)
	}

	if err = gg.SavePNG(imgPath, img); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
