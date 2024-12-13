package gpx_export

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/pkg/errors"
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
	logger     log.Logger
	log        *log.Helper
}

func NewGpxExport(logger log.Logger, app string, user, pwd string) *GpxExport {
	return &GpxExport{
		logger:   logger,
		log:      log.NewHelper(logger),
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

func (e *GpxExport) Run(isAll bool) error {
	if e.username == "" || e.password == "" {
		return fmt.Errorf("user/pwd need")
	}
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
		appExport = NewKeep(e.logger)
	case GarminCN:
		appExport = NewGarminCn(e.logger)
	case GarminCom:
		appExport = NewGarmin(e.logger)
	//
	default:
		return errors.New(fmt.Sprintf("%s app not supported", e.app))
	}

	appExport.Init(e.gitDir, e.exportDir, e.username, e.password)
	if !appExport.Auth() {
		e.log.Info("登陆失败")
		return errors.New("auth failed, skip")
	}

	e.log.Infof("auth success, start download gpx files... to %s", e.exportDir)
	err := appExport.Run(isAll)
	if err != nil {
		return err
	}
	return nil
}
