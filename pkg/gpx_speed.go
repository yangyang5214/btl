package pkg

import (
	"os"
	"path/filepath"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
)

type GpxSpeed struct {
	gpxFile string
	speed   float64
	log     *log.Helper
	result  string
}

func NewGpxSpeed(gpxFile string, speed int32, logger log.Logger) *GpxSpeed {
	return &GpxSpeed{
		gpxFile: gpxFile,
		speed:   float64(speed) / 3.6,
		log:     log.NewHelper(logger),
		result:  "result.gpx",
	}
}

func (s *GpxSpeed) SetResultFile(p string) {
	s.result = p
}

func (s *GpxSpeed) getAllPoints(gpxData *gpx.GPX) ([]gpx.GPXPoint, error) {
	if len(gpxData.Tracks) == 0 {
		return nil, errors.New("no tracks found")
	}

	track := gpxData.Tracks[0]
	if len(track.Segments) == 0 {
		return nil, errors.New("no segments found")
	}

	points := track.Segments[0]
	if len(points.Points) == 0 {
		return nil, errors.New("no points found")
	}
	return points.Points, nil
}

func (s *GpxSpeed) process() (*gpx.GPX, error) {
	p, err := filepath.Abs(s.gpxFile)
	if err != nil {
		s.log.Errorf("get file abs path err %+v", err)
		return nil, err
	}
	gpxData, err := gpx.ParseFile(p)
	if err != nil {
		s.log.Errorf("gpx ParseFile err %+v", err)
		return nil, errors.WithStack(err)
	}
	points, err := s.getAllPoints(gpxData)
	if err != nil {
		s.log.Errorf("get all points err %+v", err)
		return nil, errors.WithStack(err)
	}

	prePoint := points[0]
	startTime := prePoint.Timestamp

	resultPoints := []gpx.GPXPoint{
		{
			Point: gpx.Point{
				Latitude:  prePoint.Latitude,
				Longitude: prePoint.Longitude,
				Elevation: prePoint.Elevation,
			},
			Timestamp: startTime,
		},
	}
	for i := 1; i < len(points); i++ {
		curPoint := points[i]

		distLen := prePoint.Distance3D(&curPoint.Point)
		ts := distLen * 1000 / s.speed
		curTime := startTime.Add(time.Duration(ts) * time.Millisecond)
		resultPoints = append(resultPoints, gpx.GPXPoint{
			Point: gpx.Point{
				Latitude:  curPoint.Latitude,
				Longitude: curPoint.Longitude,
				Elevation: curPoint.Elevation,
			},
			Timestamp: curTime,
		})

		prePoint = curPoint
		startTime = curTime
	}

	gpxData.Tracks[0].Segments[0].Points = resultPoints
	return gpxData, nil
}

func (s *GpxSpeed) Run() error {
	gpxData, err := s.process()
	if err != nil {
		return errors.WithStack(err)
	}
	newXml, err := gpxData.ToXml(gpx.ToXmlParams{
		Indent: true,
	})
	return os.WriteFile(s.result, newXml, 0755)
}
