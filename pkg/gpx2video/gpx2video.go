package gpx2video

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tkrajina/gpxgo/gpx"
	"math"
	"os"
	"os/exec"
	"path"
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
	Points   []Point
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

func ParseGPX(gpxData *gpx.GPX) (*Session, error) {
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
	var avgSpeed float64
	if movingData.MovingTime != 0 {
		avgSpeed = movingData.MovingDistance / movingData.MovingTime
	}

	return &Session{
		avgSpeed: avgSpeed,
		maxSpeed: movingData.MaxSpeed,
		Points:   points,
	}, nil
}

func ParseFitBytes(mergeFitCmd string, byteData []byte, logrHelper *log.Helper) (*Session, error) {
	p := path.Join(os.TempDir(), fmt.Sprintf("%d.fit", time.Now().Unix()))
	defer func() {
		_ = os.Remove(p)
	}()
	err := os.WriteFile(p, byteData, 0755)
	if err != nil {
		return nil, err
	}
	return ParseFit(mergeFitCmd, p, logrHelper)
}

func ParseFit(mergeFitCmd, fitFile string, logrHelper *log.Helper) (*Session, error) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	gpx2fitCmd := fmt.Sprintf("%s -p %s", mergeFitCmd, fitFile)
	cmd := exec.Command("/bin/bash", "-c", gpx2fitCmd)
	cmd.Dir = tempDir

	var logBuffer bytes.Buffer
	cmd.Stdout = &logBuffer
	cmd.Stderr = &logBuffer

	logrHelper.Infof("run cmd: <%s>", gpx2fitCmd)
	err = cmd.Run()
	if err != nil {
		logrHelper.Errorf("run merge-fit err: %s", logBuffer.String())
		return nil, errors.WithStack(err)
	}
	logrHelper.Infof("parser fit point: %s", logBuffer.String())

	sessionFile := path.Join(tempDir, "record.json")
	dates, err := os.ReadFile(sessionFile)
	if err != nil {
		logrHelper.Errorf("read file <%s> err", sessionFile)
		return nil, err
	}

	result := gjson.ParseBytes(dates)

	var ps []Point

	var maxSpeed float64
	var speedFlag float64

	points := result.Array()
	for _, point := range points {
		fields := point.Get("fields").Array()

		m := make(map[string]any)
		for _, field := range fields {
			m[field.Get("name").Str] = field.Get("values").Array()[0].Int()
		}

		lat, ok := m["position_lat"]
		if !ok || lat.(int64) == 0 {
			continue
		}

		lng, ok := m["position_long"]
		if !ok || lng.(int64) == 0 {
			continue
		}

		p := Point{
			Latitude:  float64(lat.(int64)) / 11930465,
			Longitude: float64(lng.(int64)) / 11930465,
		}

		speed, ok := m["speed"]
		if ok {
			speedVal := float64(speed.(int64)) / 1_000
			p.Speed = speedVal
			if speedVal > maxSpeed {
				maxSpeed = speedVal
			}
			speedFlag = speedFlag + speedVal
		}

		ps = append(ps, p)
	}
	logrHelper.Infof("all point size %d", len(ps))
	return &Session{
		maxSpeed: maxSpeed,
		avgSpeed: speedFlag / float64(len(ps)),
		Points:   ps,
	}, nil
}

// 经纬度转换为直角坐标系（墨卡托投影）
func mercatorProjection(lat, lon float64) (float64, float64) {
	x := R * lon * math.Pi / 180
	y := R * math.Log(math.Tan(math.Pi/4+lat*math.Pi/360))
	return x, y
}

func genImageBound(session *Session) *ImageBound {
	points := session.Points
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
