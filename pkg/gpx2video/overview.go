package gpx2video

import (
	"github.com/fogleman/gg"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/yangyang5214/btl/pkg"
	"math"
)

type ImgOverview struct {
	log       *log.Helper
	width     int
	height    int
	session   *Session
	resultPng string
}

func NewImgOverview(session *Session, logger log.Logger) *ImgOverview {
	return &ImgOverview{
		log:       log.NewHelper(logger),
		session:   session,
		width:     1440,
		height:    900,
		resultPng: "result.png",
	}
}

func (s *ImgOverview) SetImgPath(resultPng string) {
	s.resultPng = resultPng
}

func (s *ImgOverview) Run() error {
	if len(s.session.Points) < 10 {
		return nil //跳过极少的点
	}
	s.log.Infof("all Points size %d", len(s.session.Points))
	imgBound := genImageBound(s.session)
	s.log.Infof("max/avg speed is %f,%f", imgBound.maxSpeed, imgBound.avgSpeed)

	imgBound.width = s.width
	imgBound.height = s.height

	err := plotImage(imgBound, s.resultPng)
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

	dc.SetLineWidth(12)

	size := len(xPoints)

	step := 1
	for i := step; i < size; i += step {
		dc.SetRGB(1, 0, 0)
		x1 := (xPoints[i-step]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
		y1 := (float64(height) - (yPoints[i-step]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2

		x2 := (xPoints[i]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
		y2 := (float64(height) - (yPoints[i]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	circleSize := 15.0

	err := pkg.LoadFontFace(dc, 16)
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
	endX := (xPoints[size-step]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
	endY := (float64(height) - (yPoints[size-step]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2
	dc.SetRGB(1, 0, 0)
	dc.DrawCircle(endX, endY, circleSize)
	dc.Fill()
	dc.SetRGB(1, 1, 1) // 设置字体颜色为白色
	dc.DrawStringAnchored("终", endX, endY, 0.5, 0.5)

	return dc.SavePNG(outputImagePath)
}
