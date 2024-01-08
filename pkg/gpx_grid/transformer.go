package gpx_grid

import (
	"fmt"
	"math"

	"github.com/icholy/utm"
	"github.com/tkrajina/gpxgo/gpx"
)

type XY struct {
	X, Y int
}

func (xy *XY) String() string {
	return fmt.Sprintf("%v,%v", xy.X, xy.Y)
}

type Transformer struct {
	xys    []*XY
	bounds *XY

	minX, minY, maxX, maxY int

	zoom    float64
	gpxData *gpx.GPX
}

func NewTransformer(gpxData *gpx.GPX) *Transformer {
	return &Transformer{
		gpxData: gpxData,
		minX:    math.MaxInt,
		minY:    math.MaxInt,
		maxX:    0,
		maxY:    0,
		zoom:    100.0,
	}
}

func (t *Transformer) Parser() {
	for _, track := range t.gpxData.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				east, north, _ := utm.ToUTM(point.Latitude, point.Longitude)
				x, y := int(north/t.zoom), int(east/t.zoom)

				t.minX = min(x, t.minX)
				t.maxX = max(x, t.maxX)

				t.minY = min(y, t.minY)
				t.maxY = max(y, t.maxY)

				t.xys = append(t.xys, &XY{
					X: x,
					Y: y,
				})
			}
		}
	}

	t.bounds = &XY{
		X: t.maxX - t.minX,
		Y: t.maxY - t.minY,
	}
}

func (t *Transformer) GetBounds() *XY {
	return t.bounds
}

func (t *Transformer) Points() []*XY {
	var r []*XY
	for _, xy := range t.xys {
		r = append(r, &XY{
			X: xy.X - t.minX,
			Y: xy.Y - t.minY,
		})
	}
	return r
}
