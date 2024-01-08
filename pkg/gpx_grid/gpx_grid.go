package gpx_grid

import (
	"sort"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
	"golang.org/x/image/colornames"
)

type GpxGrid struct {
	files         []string
	width, height int
	gpxs          []*gpx.GPX
}

func NewGpxGrid() *GpxGrid {
	return &GpxGrid{
		width:  1000,
		height: 2000,
	}
}

func (g *GpxGrid) SetFiles(files []string) {
	g.files = files
}

func (g *GpxGrid) SetGpxDatas(gpxs []*gpx.GPX) {
	g.gpxs = gpxs
}

func (g *GpxGrid) sortGpx() {
	m := make(map[int]*gpx.GPX)

	var keys []int
	for _, gpxData := range g.gpxs {
		m[gpxData.GetTrackPointsNo()] = gpxData
		keys = append(keys, gpxData.GetTrackPointsNo())
	}

	sort.Ints(keys)

	var r []*gpx.GPX
	for _, key := range keys {
		r = append(r, m[key])
	}
	g.gpxs = r
}

func (g *GpxGrid) Run() error {
	dc := gg.NewContext(g.width, g.height)

	if len(g.gpxs) == 0 {
		for _, file := range g.files {
			gpxData, err := gpx.ParseFile(file)
			if err != nil {
				return err
			}
			g.gpxs = append(g.gpxs, gpxData)
		}
	}

	g.sortGpx()

	if len(g.gpxs) == 0 {
		return errors.New("no gpx file")
	}

	log.Infof("all files count: %d", len(g.gpxs))

	var startX, startY int = 0, 0 //todo

	var maxY = 0

	for _, gpxData := range g.gpxs {
		transformer := NewTransformer(gpxData)
		transformer.Parser()

		for _, point := range transformer.Points() {
			x, y := float64(point.X+startX), float64(point.Y+startY)
			dc.DrawPoint(x, y, 2)
			dc.SetColor(colornames.Blue)
		}

		bounds := transformer.GetBounds()
		log.Infof("bounds is %s", bounds.String())

		startX = startX + bounds.X

		maxY = max(maxY, bounds.Y)

		if startX > g.width {
			startX = 0
			startY = startY + maxY
		}
	}

	dc.Fill()

	return dc.SavePNG("result.png")
}
