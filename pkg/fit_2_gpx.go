package pkg

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
	"os"
	"os/exec"
	"path"
	"time"
)

type Fit2Gpx struct {
	fitFile   string
	log       *log.Helper
	workDir   string
	runFile   string
	resultGpx string
}

//go:embed script/fit_parser.py
var fit2gpx string

func NewFit2Gpx(fitFile string, logger log.Logger) *Fit2Gpx {
	homeDir, _ := os.UserHomeDir()
	workDir := path.Join(homeDir, ".fit2gpx")
	_ = os.MkdirAll(workDir, 0755)
	return &Fit2Gpx{
		fitFile:   fitFile,
		log:       log.NewHelper(logger),
		workDir:   workDir,
		runFile:   path.Join(workDir, "fit_parser.py"),
		resultGpx: "result.gpx",
	}
}

func (s *Fit2Gpx) SetResultPath(resultPath string) {
	s.resultGpx = resultPath
}

func (s *Fit2Gpx) init() error {
	_, err := os.Stat(s.runFile)
	if err == nil {
		return nil //skip
	}
	f, err := os.Create(s.runFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fit2gpx)
	if err != nil {
		return err
	}
	cmds := []string{
		"pip3 install garmin-fit-sdk",
	}
	for _, cmd := range cmds {
		err = exec.Command("/bin/bash", "-c", cmd).Run()
		if err != nil {
			s.log.Errorf("run cmd %s failed", cmd)
			return errors.WithStack(err)
		}
	}
	return nil
}

type point struct {
	ts       int64
	lat      float64
	lng      float64
	distance float64
	speed    float64
	altitude float64
}

func (s *Fit2Gpx) process() ([]*point, error) {
	fpath := path.Join("/tmp", fmt.Sprintf("%d.json", time.Now().UnixMilli()))
	newFile, err := os.Create(fpath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = newFile.Close()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Remove(fpath)
	}()

	cmd := fmt.Sprintf("python3 %s %s %s", s.runFile, s.fitFile, fpath)
	s.log.Infof("run cmd %s", cmd)
	err = exec.Command("/bin/bash", "-c", cmd).Run()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	bytesData, err := os.ReadFile(fpath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var points []*point
	if len(bytesData) == 0 {
		return nil, errors.New("fit file parser failed")
	}
	err = json.Unmarshal(bytesData, &points)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return points, nil
}

func (s *Fit2Gpx) Run() error {
	if len(s.fitFile) == 0 {
		s.log.Infof("not fit file")
		return nil
	}
	err := s.init()
	if err != nil {
		return errors.WithStack(err)
	}
	points, err := s.process()
	if err != nil {
		return errors.WithStack(err)
	}
	gpxFile, err := gpx.ParseString("") //todo
	if err != nil {
		return errors.WithStack(err)
	}
	var gpxPoints []*gpx.GPXPoint
	for _, p := range points {
		gpxPoints = append(gpxPoints, &gpx.GPXPoint{
			Point: gpx.Point{
				Latitude:  p.lat,
				Longitude: p.lng,
				Elevation: *gpx.NewNullableFloat64(p.altitude),
			},
			Timestamp: time.UnixMilli(p.ts),
		})
	}

	newXml, err := gpxFile.ToXml(gpx.ToXmlParams{
		Indent: true,
	})
	f, err := os.Create(s.resultGpx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	_, _ = f.Write(newXml)
	return nil
}
