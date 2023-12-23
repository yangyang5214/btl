package pkg

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/yangyang5214/btl/pkg/utils"
)

type GpxZip struct {
	file       string
	expectSize int32

	outFile string
}

func NewGpxZip(file string, expectSize int32) *GpxZip {
	return &GpxZip{
		file:       file,
		expectSize: expectSize,
	}
}

func (g *GpxZip) SetOutFile(outFile string) {
	g.outFile = outFile
}

func (g *GpxZip) getOutFile() string {
	if g.outFile != "" {
		return g.outFile
	}
	return path.Join(path.Dir(g.file), "new_"+path.Base(g.file))
}

func (g *GpxZip) Run() error {
	log.Infof("start process file %s, limit size %d", g.file, g.expectSize)
	f, err := os.Stat(g.file)
	if err != nil {
		return err
	}
	fileSize := f.Size() / 1024 / 1024
	if fileSize < int64(g.expectSize) {
		log.Infof("size is not larger than expected")
		return nil
	}

	gpxDatas, err := utils.ParseGpxData([]string{
		g.file,
	})
	if len(gpxDatas) == 0 {
		return fmt.Errorf("not found gpx file")
	}
	date := g.zipPoints(gpxDatas[0], int(fileSize+1))

	resultFile, err := os.Create(g.getOutFile())
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

func (g *GpxZip) zipPoints(gpxObj *gpx.GPX, fileSize int) []byte {
	points := gpxObj.Tracks[0].Segments[0].Points

	length := len(points)
	finalPointLen := (length / fileSize) * 25
	step := length / (length - finalPointLen)
	log.Infof("original size %d, point size %d->%d,zip step is %d", fileSize, length, finalPointLen, step)

	var result []gpx.GPXPoint
	for i := 0; i < len(points); i++ {
		if i%step != 0 {
			result = append(result, points[i])
		}
	}
	gpxObj.Tracks[0].Segments[0].Points = result
	date, _ := gpxObj.ToXml(gpx.ToXmlParams{
		Indent: true, //相差不大
	})
	return date
}
