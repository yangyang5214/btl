package gpx2video

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
	"math"
	"os"
	"path"
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
	session, err := parseGPX(s.gpxFile)
	if err != nil {
		return errors.WithStack(err)
	}
	s.log.Infof("all points size %d", len(session.points))
	imgBound := genImageBound(session)
	s.log.Infof("max/avg speed is %f,%f", imgBound.maxSpeed, imgBound.avgSpeed)

	imgBound.width = s.width
	imgBound.height = s.height

	err = plotImage(imgBound, "result.png")
	if err != nil {
		return errors.WithStack(err)
	}
	s.log.Infof("gen gpx img overview success")
	return nil
}

func plotImage(imgBound *ImageBound, outputImagePath string) error {
	width, height := imgBound.width, imgBound.height
	dc := gg.NewContext(width, height)
	dc.SetRGBA(0, 0, 0, 0) // 设置背景为完全透明
	dc.Clear()

	// 计算缩放比例
	scaleX := float64(width-50) / (imgBound.maxX - imgBound.minX)
	scaleY := float64(height-50) / (imgBound.maxY - imgBound.minY)
	scale := math.Min(scaleX, scaleY)

	xPoints := imgBound.xPoints
	yPoints := imgBound.yPoints

	minX, maxX, minY, maxY := imgBound.minX, imgBound.maxX, imgBound.minY, imgBound.maxY

	dc.SetLineWidth(8)

	speeds := imgBound.speeds
	maxSpeed := imgBound.maxSpeed
	avgSpeed := imgBound.avgSpeed

	size := len(xPoints)
	for i := 1; i < size; i++ {
		speed := speeds[i-1]
		if speed < avgSpeed {
			normalizedSpeed := speed / avgSpeed
			r := normalizedSpeed
			g := 1.0
			b := 0.0
			dc.SetRGB(r, g, b) // 绿色到黄色的渐变
		} else {
			normalizedSpeed := (speed - avgSpeed) / (maxSpeed - avgSpeed)
			r := 1.0
			g := 1.0 - normalizedSpeed
			b := 0.0
			dc.SetRGB(r, g, b) // 黄色到红色的渐变
		}

		x1 := (xPoints[i-1]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
		y1 := (float64(height) - (yPoints[i-1]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2

		x2 := (xPoints[i]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
		y2 := (float64(height) - (yPoints[i]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	circleSize := 15.0

	err := loadFontFace(dc, 16)
	if err != nil {
		return errors.WithStack(err)
	}

	// 标记起点
	startX := (xPoints[0]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
	startY := (float64(height) - (yPoints[0]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2
	dc.SetRGB(0, 1, 0)
	dc.DrawCircle(startX, startY, circleSize)
	dc.Fill()
	dc.SetRGB(1, 1, 1) // 设置字体颜色为白色
	dc.DrawStringAnchored("起", startX, startY, 0.5, 0.5)

	// 标记终点
	endX := (xPoints[size-1]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
	endY := (float64(height) - (yPoints[size-1]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2
	dc.SetRGB(1, 0, 0)
	dc.DrawCircle(endX, endY, circleSize)
	dc.Fill()
	dc.SetRGB(1, 1, 1) // 设置字体颜色为白色
	dc.DrawStringAnchored("终", endX, endY, 0.5, 0.5)

	return dc.SavePNG(outputImagePath)
}

func ttfPath() string {
	homeDir, _ := os.UserHomeDir()
	return path.Join(homeDir, ".ttf", "chinese.ttf")
}

func loadFontFace(dc *gg.Context, points float64) error {
	p := ttfPath()
	_, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("ttf path %s not exist", p)
	}
	err = dc.LoadFontFace(p, points)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
