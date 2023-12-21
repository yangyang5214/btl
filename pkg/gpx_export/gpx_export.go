package gpx_export

import (
	"fmt"
	"github.com/pkg/errors"
)

type GpxExport struct {
	app       App
	exportDir string
}

func NewGpxExport(app string) *GpxExport {
	return &GpxExport{
		app: app,
	}
}

func (e *GpxExport) setExportDir(exportDir string) {
	e.exportDir = exportDir
}

func (e *GpxExport) Run() error {
	switch e.app {
	case Strava:
	//
	default:
		return errors.New(fmt.Sprintf("%s app not supported", e.app))
	}
	return nil
}
