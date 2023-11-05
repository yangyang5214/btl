package pkg

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/yangyang5214/btl/pkg/utils"
	sm "github.com/yangyang5214/go-staticmaps"
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

	weight  float64
	bgColor color.Color

	width, height int
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
		bgColor:       color.RGBA{R: 255, G: 255, B: 255, A: 255}, //white
	}
}

func (g *GpxMap) SetWeight(weight float64) {
	g.weight = weight
}

func (g *GpxMap) SetBgColor(cc color.Color) {
	g.bgColor = cc
}

func (g *GpxMap) SetSize(width, height int) {
	g.width = width
	g.height = height
}

func (g *GpxMap) getWeight(post []s2.LatLng) float64 {
	if g.weight != 0 {
		return g.weight
	}
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

func (g *GpxMap) Process() (image.Image, error) {
	gpxDatas, err := utils.ParseGpxData(g.files)
	if err != nil {
		return nil, err
	}
	positions, err := utils.ParsePositions(gpxDatas)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	//gen stat
	err = g.genStat()
	if err != nil {
		return nil, err
	}

	if g.width == 0 {
		width, height := utils.GenWidthHeight(positions)
		log.Infof("use height=%d, width=%d", height, width)
		g.width = width
		g.height = height
	}
	g.smCtx.SetSize(g.width, g.height)

	for index, post := range positions {
		g.smCtx.AddObject(sm.NewPath(post, utils.GetColor(index, g.colors), g.weight))
	}

	titleProvider, ok := g.tileProviders[g.titleName]
	if !ok {
		titleProvider = sm.NewTileProviderOpenStreetMaps()
	}

	titleProvider.Attribution = g.attribution
	titleProvider.BackGroundColor = g.bgColor

	g.smCtx.SetTileProvider(titleProvider)
	return g.smCtx.Render()
}

func (g *GpxMap) Run(imgPath string) error {
	img, err := g.Process()
	if err != nil {
		return err
	}
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
