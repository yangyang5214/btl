package gpx2video

import (
	"github.com/fogleman/gg"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
	"math"
)

type ImgOverview struct {
	gpxFile string
	log     *log.Helper
	gpxData *gpx.GPX
	width   int
	height  int
}

func NewImgOverview(gpxFile string, logger log.Logger) (*ImgOverview, error) {
	gpxData, err := gpx.ParseFile(gpxFile)
	if err != nil {
		return nil, err
	}
	return &ImgOverview{
		gpxFile: gpxFile,
		log:     log.NewHelper(logger),
		gpxData: gpxData,
		width:   800,
		height:  600,
	}, nil
}

func (s *ImgOverview) Run() error {
	s.log.Infof("gpx2video img cmd, gpx file is %s", s.gpxFile)

	//get all points
	points, err := parseGPX(s.gpxFile)
	if err != nil {
		return errors.WithStack(err)
	}
	s.log.Infof("all points size %d", len(points))

	imgBound := genImageBound(points)
	imgBound.width = s.width
	imgBound.height = s.height

	err = plotImage(imgBound, "result.png")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func plotImage(imgBound *ImageBound, outputImagePath string) error {
	width, height := imgBound.width, imgBound.height
	dc := gg.NewContext(width, height)
	dc.SetRGBA(0, 0, 0, 0) // 设置背景为完全透明
	dc.Clear()

	// 计算缩放比例
	scaleX := float64(width) / (imgBound.maxX - imgBound.minX)
	scaleY := float64(height) / (imgBound.maxY - imgBound.minY)
	scale := math.Min(scaleX, scaleY)

	xPoints := imgBound.xPoints
	yPoints := imgBound.yPoints

	minX, maxX, minY, maxY := imgBound.minX, imgBound.maxX, imgBound.minY, imgBound.maxY

	dc.SetRGB(1, 0, 0)

	size := len(xPoints)
	for i := 1; i < size; i++ {
		x1 := (xPoints[i-1]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
		y1 := (float64(height) - (yPoints[i-1]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2

		x2 := (xPoints[i]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
		y2 := (float64(height) - (yPoints[i]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2
		dc.DrawLine(x1, y1, x2, y2)
	}
	dc.Stroke()
	return dc.SavePNG(outputImagePath)
}
