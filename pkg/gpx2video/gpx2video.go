package gpx2video

import (
	_ "embed"
	"github.com/tkrajina/gpxgo/gpx"
	"math"
	"os"
	"time"
)

const R = 6378137

type Point struct {
	Latitude  float64
	Longitude float64
	Elevation float64

	Timestamp time.Time

	Speed float64
}

type Session struct {
	points   []Point
	maxSpeed float64
	avgSpeed float64
}

type ImageBound struct {
	xPoints, yPoints       []float64
	minX, maxX, minY, maxY float64
	width, height          int

	speeds []float64

	maxSpeed float64
	avgSpeed float64
}

// 解析 GPX 文件
func parseGPX(filePath string) (*Session, error) {
	gpxFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer gpxFile.Close()

	gpxData, err := gpx.Parse(gpxFile)
	if err != nil {
		return nil, err
	}

	var points []Point
	for _, track := range gpxData.Tracks {
		for _, segment := range track.Segments {
			for index, p := range segment.Points {
				points = append(points, Point{
					Latitude:  p.Latitude,
					Longitude: p.Longitude,
					Timestamp: p.Timestamp,
					Elevation: p.Elevation.Value(),
					Speed:     segment.Speed(index),
				})
			}
		}
	}

	movingData := gpxData.MovingData()
	return &Session{
		avgSpeed: movingData.MovingDistance / movingData.MovingTime,
		maxSpeed: movingData.MaxSpeed,
		points:   points,
	}, nil
}

// 经纬度转换为直角坐标系（墨卡托投影）
func mercatorProjection(lat, lon float64) (float64, float64) {
	x := R * lon * math.Pi / 180
	y := R * math.Log(math.Tan(math.Pi/4+lat*math.Pi/360))
	return x, y
}

func genImageBound(session *Session) *ImageBound {
	points := session.points
	var (
		xPoints, yPoints, speeds []float64
	)
	for _, point := range points {
		x, y := mercatorProjection(point.Latitude, point.Longitude)
		xPoints = append(xPoints, x)
		yPoints = append(yPoints, y)

		speeds = append(speeds, point.Speed)
	}

	// 计算坐标范围
	minX, maxX := xPoints[0], xPoints[0]
	minY, maxY := yPoints[0], yPoints[0]
	for i := range xPoints {
		if xPoints[i] < minX {
			minX = xPoints[i]
		}
		if xPoints[i] > maxX {
			maxX = xPoints[i]
		}
		if yPoints[i] < minY {
			minY = yPoints[i]
		}
		if yPoints[i] > maxY {
			maxY = yPoints[i]
		}
	}

	return &ImageBound{
		xPoints:  xPoints,
		yPoints:  yPoints,
		minX:     minX,
		maxX:     maxX,
		minY:     minY,
		maxY:     maxY,
		speeds:   speeds,
		maxSpeed: session.maxSpeed,
		avgSpeed: session.avgSpeed,
	}
}
