package pkg

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/yangyang5214/btl/pkg/utils"
)

type GpxZip struct {
	file string
	step int

	outFile string

	log *log.Helper
}

func NewGpxZip(file string, step int) *GpxZip {
	if step == 0 {
		step = 2 //default 2
	}
	return &GpxZip{
		file: file,
		step: step,
		log:  log.NewHelper(log.DefaultLogger),
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
	g.log.Infof("start process file %s, use step %d", g.file, g.step)
	gpxDatas, err := utils.ParseGpxData([]string{
		g.file,
	})
	if len(gpxDatas) == 0 {
		return fmt.Errorf("not found gpx file")
	}
	date := g.zipPoints(gpxDatas[0])

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

func (g *GpxZip) zipPoints(gpxObj *gpx.GPX) []byte {
	points := gpxObj.Tracks[0].Segments[0].Points
	var result []gpx.GPXPoint
	for i := 0; i < len(points); i++ {
		if i%g.step != 0 {
			result = append(result, points[i])
		}
	}
	gpxObj.Tracks[0].Segments[0].Points = result
	date, _ := gpxObj.ToXml(gpx.ToXmlParams{
		Indent: true, //相差不大
	})
	return date
}
