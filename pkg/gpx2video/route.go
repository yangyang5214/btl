package gpx2video

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
	"math"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type RouteVideo struct {
	log           *log.Helper
	workDir       string
	width, height int
	outMp4        string
	gpxData       *gpx.GPX
}

func NewRouteVideo(gpxData *gpx.GPX, logger log.Logger, workDir string) *RouteVideo {
	return &RouteVideo{
		gpxData: gpxData,
		workDir: workDir,
		log:     log.NewHelper(logger),
		width:   800,
		height:  600,
		outMp4:  "out.mp4",
	}
}

func (s *RouteVideo) Run() error {
	session, err := ParseGPX(s.gpxData)
	if err != nil {
		return errors.WithStack(err)
	}
	points := session.points
	s.log.Infof("all points size %d", len(points))

	imgBound := genImageBound(session)

	prePoint := points[0]
	err = s.genImage(0, points, imgBound)
	if err != nil {
		return errors.WithStack(err)
	}

	for i := 1; i < len(points); i++ {
		curPoint := points[i]

		err = s.genImage(i, points, imgBound)
		if err != nil {
			return errors.WithStack(err)
		}

		subTs := curPoint.Timestamp.Sub(prePoint.Timestamp)
		subSeconds := int(subTs.Seconds())
		err = s.copyPre(subSeconds, prePoint.Timestamp)
		if err != nil {
			return errors.WithStack(err)
		}

		prePoint = curPoint
	}

	err = s.genVideo(s.workDir, s.outMp4)
	if err != nil {
		return err
	}
	return nil
}

func (s *RouteVideo) genImage(index int, points []Point, imgBound *ImageBound) error {
	imgPath := s.genImagePath(points[index].Timestamp)
	err := s.plotImage(imgBound, index, imgPath)
	if err != nil {
		return err
	}
	return nil
}

func (s *RouteVideo) plotImage(imgBound *ImageBound, flag int, outputImagePath string) error {
	width, height := s.width, s.height
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

	dc.SetRGB(1, 0, 0) // 设置轨迹点颜色为红色
	for i := range xPoints {
		// 将坐标转换到图像空间
		x := (xPoints[i]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
		y := (float64(height) - (yPoints[i]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2
		dc.DrawPoint(x, y, 1)
		dc.Fill()

		if i == flag {
			dc.SetRGB(1, 1, 0)
		}
	}
	return dc.SavePNG(outputImagePath)
}

func (s *RouteVideo) genImagePath(ts time.Time) string {
	return path.Join(s.workDir, fmt.Sprintf("%d.png", ts.Unix()))
}

func (s *RouteVideo) copyPre(num int, startTime time.Time) error {
	prePath := s.genImagePath(startTime)
	for i := 1; i < num; i++ {
		curPath := s.genImagePath(startTime.Add(time.Duration(i) * time.Second))
		cpCmd := fmt.Sprintf("cp %s %s", prePath, curPath)
		err := exec.Command("/bin/bash", "-c", cpCmd).Run()
		if err != nil {
			s.log.Errorf("run cmd <%s> failed", cpCmd)
			return err
		}
	}
	return nil
}

func (s *RouteVideo) genVideo(workDir string, videoPath string) error {
	s.log.Infof("gen mp4 start")
	f, err := os.CreateTemp("", "*.txt")
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = os.Remove(f.Name())
	}()
	entries, err := os.ReadDir(workDir)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".png") {
			_, _ = f.WriteString(fmt.Sprintf("file '%s'\n", path.Join(workDir, entry.Name())))
			_, _ = f.WriteString("duration 1\n")
		}
	}

	ffmepgCmd := fmt.Sprintf("ffmpeg -f concat -safe 0 -i %s  %s", f.Name(), videoPath)
	s.log.Infof("run cmd: %s", ffmepgCmd)
	err = exec.Command("/bin/bash", "-c", ffmepgCmd).Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
