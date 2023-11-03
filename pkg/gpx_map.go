package pkg

import (
	"fmt"
	"image"
	"image/color"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/yangyang5214/btl/pkg/utils"
	"golang.org/x/image/font/gofont/goregular"
)

type Stat struct {
	distance     float64
	timeOfSecond float64
	count        int
}

func (s *Stat) String() string {
	return fmt.Sprintf("distance %vkm, time %vh, count %d", int(s.distance)/1000, int(s.timeOfSecond)/3600, s.count)
}

type GpxMap struct {
	files         []string
	gpxData       []*gpx.GPX
	smCtx         *sm.Context
	attribution   string
	titleName     string
	tileProviders map[string]*sm.TileProvider
	colors        []color.Color
	stat          *Stat
}

func NewGpxMap(files []string, attribution, titleName string, colors []color.Color) *GpxMap {
	if len(colors) == 0 {
		defaultColor, _ := sm.ParseColorString("green")
		colors = []color.Color{defaultColor}
	}
	return &GpxMap{
		files:         files,
		smCtx:         sm.NewContext(),
		attribution:   attribution,
		titleName:     titleName,
		tileProviders: sm.GetTileProviders(),
		colors:        colors,
	}
}

func (g *GpxMap) getWeight(post []s2.LatLng) float64 {
	var weight float64
	defer func() {
		log.Infof("pointCount size is %d, line weight %v", len(post), weight)
	}()
	size := len(post)
	weight = float64(size/10_000) + 2.0
	return weight
}

func (g *GpxMap) genStat() error {
	stat := &Stat{}
	for _, gd := range g.gpxData {
		md := gd.MovingData()
		stat.distance = stat.distance + md.MovingDistance + md.StoppedDistance
		stat.timeOfSecond = stat.timeOfSecond + md.MovingTime + md.StoppedTime
		gpx.GetGpxElementInfo("", gd)
	}
	stat.count = len(g.gpxData)
	g.stat = stat
	return nil
}

func (g *GpxMap) getColor(index int) color.Color {
	size := len(g.colors)
	return g.colors[index%size]
}

func (g *GpxMap) Run(imgPath string) error {
	gpxDatas, err := utils.ParseGpxData(g.files)
	if err != nil {
		return err
	}
	positions, err := utils.ParsePositions(gpxDatas)
	if err != nil {
		return errors.WithStack(err)
	}
	//gen stat
	err = g.genStat()
	if err != nil {
		return err
	}

	width, height := utils.GenWidthHeight(positions)
	log.Infof("use height=%d, width=%d", height, width)
	g.smCtx.SetSize(width, height)

	for index, post := range positions {
		weight := g.getWeight(post)
		if height <= 1000 {
			weight = 3
		}
		g.smCtx.AddObject(sm.NewPath(post, g.getColor(index), weight))
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

	//img = g.addStat(img)
	if err = gg.SavePNG(imgPath, img); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (g *GpxMap) addStat(img image.Image) image.Image {
	log.Infof(g.stat.String())
	img = addLabel(img, 50, 100, g.stat.String())
	return img
}

func addLabel(img image.Image, x, y int, label string) image.Image {
	dc := gg.NewContextForImage(img)
	dc.SetColor(&color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	})
	font, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(font, &truetype.Options{Size: 50})
	dc.SetFontFace(face)
	dc.DrawStringAnchored(label, float64(x), float64(y), 0, 0)
	return dc.Image()
}
