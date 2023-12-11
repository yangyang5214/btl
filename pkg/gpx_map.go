package pkg

import (
	"fmt"
	"golang.org/x/image/colornames"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
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

func (s *Stat) Distance() string {
	return fmt.Sprintf("%v km", int(s.distance)/1000)
}

type GpxMap struct {
	files         []string
	smCtx         *sm.Context
	attribution   string
	titleName     string
	tileProviders map[string]*sm.TileProvider
	colors        []color.Color
	stat          *Stat

	bgColor color.Color
	weight  float64
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
		weight:        3,
	}
}

func (g *GpxMap) SetWeight(weight float64) {
	g.weight = weight
}

func (g *GpxMap) SetBgColor(cc color.Color) {
	g.bgColor = cc
}

func (g *GpxMap) genStat(gpxs []*gpx.GPX) error {
	stat := &Stat{}
	for _, gd := range gpxs {
		md := gd.MovingData()
		stat.distance = stat.distance + md.MovingDistance + md.StoppedDistance
		stat.timeOfSecond = stat.timeOfSecond + md.MovingTime + md.StoppedTime
		gpx.GetGpxElementInfo("", gd)
	}
	stat.count = len(gpxs)
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
	err = g.genStat(gpxDatas)
	if err != nil {
		return nil, err
	}

	w, h := utils.GenWidthHeight(positions)
	g.smCtx.SetSize(w, h)

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

	img = g.addStat(img)
	if err = gg.SavePNG(imgPath, img); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (g *GpxMap) addStat(img image.Image) image.Image {
	log.Infof("stat info is %s", g.stat.String())
	img = addLabel(img, 50, 50, g.stat.Distance())
	return img
}

func addLabel(img image.Image, x, y int, label string) image.Image {
	dc := gg.NewContextForImage(img)
	dc.SetColor(colornames.White)
	font, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(font, &truetype.Options{Size: 50})
	dc.SetFontFace(face)
	dc.DrawStringAnchored(label, float64(x), float64(y), 0, 0)
	return dc.Image()
}
