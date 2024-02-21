package gpx_export

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg"
)

type GpxExport struct {
	app       App
	exportDir string

	username   string
	password   string
	skipUpdate bool
	gitDir     string
	repoURL    string
}

func NewGpxExport(app string, user, pwd string) *GpxExport {
	return &GpxExport{
		app:      app,
		username: user,
		password: pwd,
		repoURL:  "git@github.com:yangyang5214/gpx_export.git",
		gitDir:   DefaultGitDir(),
	}
}

func (e *GpxExport) SetExportDir(exportDir string) {
	e.exportDir = exportDir
}

func (e *GpxExport) SkipUpdate() {
	e.skipUpdate = true
}

func (e *GpxExport) Run() error {
	if !pkg.FileExists(e.exportDir) {
		return errors.New(fmt.Sprintf("exportDir: <%s> is not exists", e.exportDir))
	}

	if !e.skipUpdate {
		err := GitClone(e.gitDir, e.repoURL)
		if err != nil {
			return err
		}
	}
	var appExport AppExport
	switch e.app {
	case Keep:
		appExport = NewKeep()
	case GarminCN:
		appExport = NewGarminCn()
	//
	default:
		return errors.New(fmt.Sprintf("%s app not supported", e.app))
	}

	appExport.Init(e.gitDir, e.exportDir, e.username, e.password)
	if !appExport.Auth() {
		log.Info("登陆失败")
		return errors.New("auth failed, skip")
	}

	log.Infof("auth success, start download gpx files... to %s", e.exportDir)
	err := appExport.Run()
	if err != nil {
		return err
	}
	return nil
}
