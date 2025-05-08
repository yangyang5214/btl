package pkg

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/yangyang5214/btl/pkg/utils"
)

type GpxMerge struct {
	currentDir string
	gpxDatas   []*gpx.GPX
	resultPath string
	log        *log.Helper
}

func NewGpxMerge(workDir string, logger log.Logger) *GpxMerge {
	return &GpxMerge{
		currentDir: workDir,
		resultPath: path.Join(workDir, "result.gpx"),
		log:        log.NewHelper(logger),
	}
}

func (g *GpxMerge) SetResultPath(p string) {
	g.log.Infof("Set result path to: %s", p)
	g.resultPath = p
}

func (g *GpxMerge) SetGpxDatas(gpxDatas []*gpx.GPX) {
	g.gpxDatas = gpxDatas
}

func (g *GpxMerge) loadGpxData() error {
	dirs, err := os.ReadDir(g.currentDir)
	if err != nil {
		return errors.WithStack(err)
	}

	var gpxFiles []string
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}
		if strings.HasSuffix(dir.Name(), ".gpx") {
			g.log.Infof("find gpx file %s", dir.Name())
			gpxFiles = append(gpxFiles, path.Join(g.currentDir, dir.Name()))
		}
	}

	if len(gpxFiles) == 0 {
		g.log.Info("not find gpx files in current directory")
		return nil
	}
	gpxDatas, err := utils.ParseGpxData(gpxFiles)
	if err != nil {
		return err
	}
	g.gpxDatas = gpxDatas
	return nil
}

func (g *GpxMerge) Run(removeFirstPoint bool) error {
	if len(g.gpxDatas) == 0 {
		err := g.loadGpxData()
		if err != nil {
			return err
		}
	}

	if len(g.gpxDatas) == 0 {
		return fmt.Errorf("no gpx files found")
	}

	g.gpxDatas = utils.SortGpx(g.gpxDatas)

	g.log.Infof("merge gpx files count: %d", len(g.gpxDatas))
	if len(g.gpxDatas) == 0 {
		return errors.New("gpx files is zero")
	}

	firstGpx := g.gpxDatas[0]
	points := firstGpx.Tracks[0].Segments[0].Points

	//追加 points
	for i := 1; i < len(g.gpxDatas); i++ {
		curPoints := g.gpxDatas[i].Tracks[0].Segments[0].Points

		g.log.Infof("Append new points count %d, for file index %d", len(curPoints), i)

		points = append(points, curPoints...)
	}

	g.log.Infof("All points count %d", len(points))

	if removeFirstPoint {
		points = points[1:]
	}

	firstGpx.Tracks[0].Segments[0].Points = points

	date, err := firstGpx.ToXml(gpx.ToXmlParams{
		Indent: true,
	})
	if err != nil {
		return err
	}

	g.log.Infof("save gpx result to: %s", g.resultPath)

	resultFile, err := os.Create(g.resultPath)
	defer resultFile.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = resultFile.Write(date)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (g *GpxMerge) sortByDate() error {
	dateMap := make(map[int64]*gpx.GPX)
	keys := make([]int64, 0, len(dateMap))
	for _, f := range g.gpxDatas {
		startTime := f.TimeBounds().StartTime.UnixMilli()
		keys = append(keys, startTime)
		dateMap[startTime] = f
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	var result []*gpx.GPX
	for index, k := range keys {
		startTime := time.UnixMilli(k).Format("2006-01-02 15:04:05")
		g.log.Infof("index %d, time: %s", index+1, startTime)
		result = append(result, dateMap[k])
	}
	g.gpxDatas = result
	return nil
}
