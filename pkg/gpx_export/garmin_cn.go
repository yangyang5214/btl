package gpx_export

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/pkg/errors"
)

// GarminCn 国区
type GarminCn struct {
	gitDir string

	username  string
	password  string
	exportDir string
	log       *log.Helper
}

func NewGarminCn(logger log.Logger) *GarminCn {
	return &GarminCn{
		log: log.NewHelper(log.With(logger, "app", "garmin_cn")),
	}
}

func (g *GarminCn) Init(gitDir string, exportDir string, username, password string) {
	g.gitDir = gitDir
	g.exportDir = exportDir
	g.username = username
	g.password = password
}

func (g *GarminCn) Auth() bool {
	cmdStr := fmt.Sprintf("python3 %s/garmin_secret.py -u '%s' -p '%s' --cn", g.gitDir, g.username, g.password)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	g.log.Infof("start run cmd: %s", cmdStr)

	out, err := cmd.Output()
	if err != nil {
		g.log.Errorf("run cmd: %s", err)
		return false
	}
	g.log.Infof("run garmin_secret out: %s", out)
	return strings.Contains(string(out), "success")
}

func (g *GarminCn) Run() error {
	cmdStr := fmt.Sprintf("python3 %s/garmin_export.py --is-cn -u '%s' -p '%s' --out %s", g.gitDir, g.username, g.password, g.exportDir)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	g.log.Infof("satrt run cmd: %s", cmdStr)

	garminExportOut, err := cmd.Output()
	if err != nil {
		g.log.Infof("run cmd: %s", err)
		return err
	}
	result := string(garminExportOut)
	result = strings.Trim(result, "")
	g.log.Infof("run garmin_export result:\n %s", result)
	if strings.HasSuffix(result, "seconds") {
		return errors.New(result)
	}
	return nil
}
