package pkg

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/yangyang5214/btl/pkg/fit2gpx"
	"os"
	"time"
)

type Fit2Kml struct {
	input  string
	logger log.Logger
	log    *log.Helper
	output string
}

func NewFit2Kml(input string, output string) *Fit2Kml {
	logger := log.DefaultLogger

	if output == "" {
		output = "result.kml"
	}
	return &Fit2Kml{
		input:  input,
		logger: logger,
		log:    log.NewHelper(logger),
		output: output,
	}
}

func (s *Fit2Kml) Run() error {
	f2g := fit2gpx.NewFit2Gpx(s.input, s.logger)
	gpxData, err := f2g.ParseToGpx()
	if err != nil {
		return err
	}

	gpxPath := fmt.Sprintf("%d.gpx", time.Now().UnixMilli())
	err = f2g.SaveGpx(gpxData, gpxPath)
	if err != nil {
		return err
	}
	defer os.Remove(gpxPath)

	opts := WithResultFile(s.output)
	err = NewGpx2Kml(gpxPath, s.logger, opts).Run()
	if err != nil {
		s.log.Errorf("Run NewGpx2Kml err: %s", err)
		return err
	}
	s.log.Infof("Run NewGpx2Kml success. save to %s", s.output)
	return nil
}
