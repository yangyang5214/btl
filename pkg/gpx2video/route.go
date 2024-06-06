package gpx2video

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tkrajina/gpxgo/gpx"
	"time"
)

type GpxVideo struct {
	gpxFile string
	log     *log.Helper
}

func NewGpxVideo(gpxFile string, logger log.Logger) *GpxVideo {
	return &GpxVideo{
		gpxFile: gpxFile,
		log:     log.NewHelper(logger),
	}
}

func (s *GpxVideo) Run() error {
	s.log.Infof("gpx2video gpx file is %s", s.gpxFile)

	points, err := s.parserPoints()
	if err != nil {
		return err
	}
	s.log.Infof("all points size %d", len(points))

	prePoint := points[0]
	for i := 1; i < len(points); i++ {
		curPoint := points[i]

		subTs := curPoint.Timestamp.Sub(prePoint.Timestamp)
		subSeconds := int(subTs.Seconds())
		s.copyPre(subSeconds, prePoint.Timestamp)

		//trans
		prePoint = curPoint
	}

	return nil
}

func (s *GpxVideo) parserPoints() ([]gpx.GPXPoint, error) {
	gpxData, err := gpx.ParseFile(s.gpxFile)
	if err != nil {
		return nil, err
	}
	return gpxData.Tracks[0].Segments[0].Points, nil
}

// copyPre 复制前置图片
func (s *GpxVideo) copyPre(num int, startTime time.Time) {
	for i := 1; i < num; i++ {
		startTime.Add(time.Duration(i) * time.Second)
	}
}
