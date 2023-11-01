package pkg

import (
	"fmt"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg/utils"
	"image"
	"image/color"
)

type GpxVideo struct {
	files         []string
	col           color.Color
	titleProvider *sm.TileProvider
}

func NewGpxVideo(files []string) *GpxVideo {
	col, _ := sm.ParseColorString("green")
	titleProvider := sm.NewTileProviderCartoDark()
	titleProvider.Attribution = ""
	return &GpxVideo{
		files:         files,
		col:           col,
		titleProvider: titleProvider,
	}
}

func (g *GpxVideo) genCenter(positions [][]s2.LatLng) (int, s2.LatLng, error) {
	smCtx := sm.NewContext()
	for _, position := range positions {
		smCtx.AddObject(sm.NewPath(position, g.col, 2))
	}
	return smCtx.DetermineZoomCenter()
}

func (g *GpxVideo) Run() error {
	gpxDatas, err := utils.ParseGpxData(g.files)
	if err != nil {
		return err
	}
	positions, err := utils.ParsePositions(gpxDatas)
	if err != nil {
		return errors.WithStack(err)
	}
	zoom, center, err := g.genCenter(positions)
	if err != nil {
		return errors.WithStack(err)
	}

	width, height := utils.GenWidthHeight(positions)
	log.Infof("use height=%d, width=%d", height, width)

	var img image.Image
	var index int64
	for _, position := range positions {
		for i := 1; i < len(position); i += 5 {
			index = index + 1
			smCtx := sm.NewContext()
			smCtx.SetSize(width, height)
			smCtx.SetTileProvider(g.titleProvider)
			smCtx.SetZoom(zoom)
			smCtx.SetCenter(center)

			smCtx.AddObject(sm.NewPath(position[0:i], g.col, 2))

			img, err = smCtx.Render()
			if err != nil {
				return errors.WithStack(err)
			}

			//img = g.addStat(img)
			if err = gg.SavePNG(fmt.Sprintf("/tmp/2/%d.png", index), img); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	//ffmpeg -framerate 30 -i %d.png -c:v libx264  out.mp4 -y
	return nil
}
