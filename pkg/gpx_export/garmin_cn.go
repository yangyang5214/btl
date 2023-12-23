package gpx_export

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg"
)

type GarminCn struct {
	repoURL    string
	runningDir string

	username  string
	password  string
	exportDir string
}

func NewGarminCn(exportDir string) *GarminCn {
	userHome, _ := os.UserHomeDir()
	return &GarminCn{
		repoURL:    "git@github.com:yangyang5214/gpx_export.git",
		runningDir: path.Join(userHome, ".gpx_export"),
		exportDir:  exportDir,
	}
}

func (g *GarminCn) checkPythonVersion() error {
	cmd := exec.Command("python3", "-c", "import sys; print('.'.join(map(str, sys.version_info[:2])))")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	versionStr := strings.TrimSpace(string(output))

	log.Infof("python version is %s", versionStr)

	// Parse major and minor version
	var major, minor int
	_, err = fmt.Sscanf(versionStr, "%d.%d", &major, &minor)
	if err != nil {
		return err
	}

	// Check if the version is greater than 3.8
	if major < 3 || (major == 3 && minor < 8) {
		return errors.New("Error: Python 3.8 or higher is required.")
	}
	return nil
}

func (g *GarminCn) cloneScripts() error {
	err := g.checkPythonVersion()
	if err != nil {
		return err
	}
	if _, err := os.Stat(g.runningDir); os.IsNotExist(err) {
		cmdStr := fmt.Sprintf("git clone %s %s", g.repoURL, g.runningDir)
		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		if _, err := cmd.Output(); err != nil {
			log.Errorf("run cmd: %s, err: %v", cmdStr, err)
			return err
		}
	} else {
		cmdStr := fmt.Sprintf("cd %s && git pull", g.runningDir)
		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		if _, err := cmd.Output(); err != nil {
			log.Errorf("run cmd: %s, err: %v", cmdStr, err)
			return err
		}
	}

	log.Info("start run pip3 install...")
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && pip3 install -r requirements.txt", g.runningDir))
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}
func (g *GarminCn) Auth(username, password string) {
	g.username = username
	g.password = password

}
func (g *GarminCn) Run() error {
	if !pkg.FileExists(g.exportDir) {
		return errors.New(fmt.Sprintf("exportDir: <%s> is not exists", g.exportDir))
	}
	err := g.cloneScripts()
	if err != nil {
		log.Errorf("clone script failed: %v", err)
		return err
	}

	cmdStr := fmt.Sprintf("python3 %s/garmin_export.py --is-cn -u '%s' -p '%s' --out %s", g.runningDir, g.username, g.password, g.exportDir)
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
