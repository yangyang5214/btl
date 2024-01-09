package pkg

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/yangyang5214/btl/pkg/utils"

	log "github.com/sirupsen/logrus"
)

type GpxMerge struct {
	currentDir string
	gpxDatas   []*gpx.GPX
	resultPath string
}

func NewGpxMerge(workDir string) *GpxMerge {
	return &GpxMerge{
		currentDir: workDir,
		resultPath: path.Join(workDir, "result.gpx"),
	}
}

func (g *GpxMerge) SetResultPath(p string) {
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
			log.Infof("find gpx file %s", dir.Name())
			gpxFiles = append(gpxFiles, path.Join(g.currentDir, dir.Name()))
		}
	}

	if len(gpxFiles) == 0 {
		log.Info("not find gpx files in current directory")
		return nil
	}
	gpxDatas, err := utils.ParseGpxData(gpxFiles)
	if err != nil {
		return err
	}
	g.gpxDatas = gpxDatas
	return nil
}

func (g *GpxMerge) Run() error {
	if len(g.gpxDatas) == 0 {
		err := g.loadGpxData()
		if err != nil {
			return err
		}
	}

	if len(g.gpxDatas) == 0 {
		return fmt.Errorf("no gpx files found")
	}

	log.Infof("merge gpx files count: %d", len(g.gpxDatas))
	gpxDatas, err := g.sortByDate(g.gpxDatas)
	if err != nil {
		return errors.WithStack(err)
	}

	firstGpx := gpxDatas[0]
	for i := 1; i < len(gpxDatas); i++ {
		currentGpx := gpxDatas[i]
		firstGpx.Tracks = append(firstGpx.Tracks, currentGpx.Tracks...)
	}

	date, err := firstGpx.ToXml(gpx.ToXmlParams{
		Indent: true, //相差不大
	})
	if err != nil {
		return err
	}

	log.Infof("save gpx result to: %s", g.resultPath)

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

func (g *GpxMerge) sortByDate(files []*gpx.GPX) ([]*gpx.GPX, error) {
	dateMap := make(map[int64]*gpx.GPX)
	for _, f := range files {
		startTime := utils.GetStartTime(f)
		dateMap[startTime] = f
	}

	//get sorted files
	keys := make([]int64, 0, len(dateMap))
	for key := range dateMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	var result []*gpx.GPX
	for _, k := range keys {
		result = append(result, dateMap[k])
	}
	return result, nil
}
