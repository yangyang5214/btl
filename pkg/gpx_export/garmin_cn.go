package gpx_export

import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
)

type GarminCn struct {
	repoURL   string
	targetDir string

	username string
	password string
	outDir   string
}

func NewGarminCn() *GarminCn {
	userHome, _ := os.UserHomeDir()
	return &GarminCn{
		repoURL:   "git@github.com:yangyang5214/gpx_export.git",
		targetDir: path.Join(userHome, ".gpx_export"),
		outDir:    "",
	}
}

func (g *GarminCn) runCmd(cmd *exec.Cmd) ([]byte, error) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return cmd.Output()
}

func (g *GarminCn) SetOutDir(outDir string) {
	g.outDir = outDir
}

func (g *GarminCn) checkPythonVersion() error {
	cmd := exec.Command("python", "-c", "import sys; print('.'.join(map(str, sys.version_info[:2])))")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	versionStr := strings.TrimSpace(string(output))

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
	if _, err := os.Stat(g.targetDir); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", g.repoURL, g.targetDir)
		if _, err := g.runCmd(cmd); err != nil {
			return err
		}
	} else {
		cmd := exec.Command("git", "-C", g.targetDir, "pull")
		if _, err := g.runCmd(cmd); err != nil {
			return err
		}
	}

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && pip3 install -r requirements.txt", g.targetDir))
	if _, err := g.runCmd(cmd); err != nil {
		return err
	}
	return nil
}
func (g *GarminCn) Auth(username, password string) {
	g.username = username
	g.password = password

}
func (g *GarminCn) Run() error {
	err := g.cloneScripts()
	if err != nil {
		log.Errorf("clone script failed: %v", err)
		return err
	}

	cmd := exec.Command(fmt.Sprintf("cd %s && python3 garmin_export.py -u %s -p %s --out %s", g.targetDir, g.username, g.password, g.outDir))

	garminExportOut, err := g.runCmd(cmd)
	if err != nil {
		return err
	}
	log.Infof("run garmin_export result: %s", string(garminExportOut))

	successFile := path.Join(g.outDir, "success")
	if _, err := os.Stat(successFile); errors.Is(err, os.ErrNotExist) {
		return errors.New(string(garminExportOut))
	}
	return nil
}
