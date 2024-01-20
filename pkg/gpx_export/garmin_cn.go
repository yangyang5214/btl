package gpx_export

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type GarminCn struct {
	gitDir string

	username  string
	password  string
	exportDir string
}

func NewGarminCn() *GarminCn {
	return &GarminCn{}
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
	log.Infof("satrt run cmd: %s", cmdStr)

	out, err := cmd.Output()
	if err != nil {
		log.Errorf("run cmd: %s", err)
		return false
	}
	log.Infof("run garmin_secret out: %s", out)
	return strings.Contains(string(out), "success")
}

func (g *GarminCn) Run() error {
	cmdStr := fmt.Sprintf("python3 %s/garmin_export.py --is-cn -u '%s' -p '%s' --out %s", g.gitDir, g.username, g.password, g.exportDir)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	log.Infof("satrt run cmd: %s", cmdStr)

	garminExportOut, err := cmd.Output()
	if err != nil {
		log.Infof("run cmd: %s", err)
		return err
	}
	result := string(garminExportOut)
	result = strings.Trim(result, "")
	log.Infof("run garmin_export result:\n %s", result)
	if strings.HasSuffix(result, "seconds") {
		return errors.New(result)
	}
	return nil
}
