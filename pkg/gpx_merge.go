package pkg

import (
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
}

type GpxFile struct {
	start   []string
	end     []string
	content []string
}

func NewGpxMerge(workDir string) *GpxMerge {
	return &GpxMerge{
		currentDir: workDir,
	}
}

func (g *GpxMerge) Run() error {
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

	gpxDatas, err = g.sortByDate(gpxDatas)
	if err != nil {
		return errors.WithStack(err)
	}

	log.Infof("merge gpx files count: %d", len(gpxFiles))

	firstGpx := gpxDatas[0]

	points := firstGpx.Tracks[0].Segments[0].Points
	for i := 1; i < len(gpxDatas); i++ {
		points = append(points, gpxDatas[i].Tracks[0].Segments[0].Points...)
	}

	firstGpx.Tracks[0].Segments[0].Points = points

	date, err := firstGpx.ToXml(gpx.ToXmlParams{
		Indent: true, //相差不大
	})
	if err != nil {
		return err
	}

	resultFile, err := os.Create(path.Join(g.currentDir, "result.gpx"))
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
