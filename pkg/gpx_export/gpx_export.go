package gpx_export

import (
	"fmt"
	"github.com/pkg/errors"
)

type GpxExport struct {
	app       App
	exportDir string

	username string
	password string
}

func NewGpxExport(app string, user, pwd string) *GpxExport {
	return &GpxExport{
		app:      app,
		username: user,
		password: pwd,
	}
}

func (e *GpxExport) SetExportDir(exportDir string) {
	e.exportDir = exportDir
}

func (e *GpxExport) Run() error {
	var appExport AppExport
	switch e.app {
	case Strava:
	//
	case GarminCN:
		appExport = NewGarminCn(e.exportDir)
	//
	default:
		return errors.New(fmt.Sprintf("%s app not supported", e.app))
	}

	appExport.Auth(e.username, e.password)
	err := appExport.Run()
	if err != nil {
		return err
	}
	return nil
}
