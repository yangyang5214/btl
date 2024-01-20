package gpx_export

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
)

type App = string

var (
	Apps = []App{
		Strava,
		GarminCN,
		Keep,
	}
)

const (
	Strava   App = "strava"
	GarminCN App = "garmin_cn"
	Keep     App = "keep"
)

type AppExport interface {
	Init(gitDir string, exportDir string, username, password string)
	Run() error
	Auth() bool
}

func checkPythonVersion() error {
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

func DefaultGitDir() string {
	userHome, _ := os.UserHomeDir()
	return path.Join(userHome, ".gpx_export")
}

func GitClone(gitDir string, gitRepoUrl string) error {
	err := checkPythonVersion()
	if err != nil {
		return err
	}
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		cmdStr := fmt.Sprintf("git clone %s %s", gitRepoUrl, gitDir)
		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		if _, err := cmd.Output(); err != nil {
			log.Errorf("run cmd: %s, err: %v", cmdStr, err)
			return err
		}
	} else {
		cmdStr := fmt.Sprintf("cd %s && git pull", gitDir)
		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		if _, err := cmd.Output(); err != nil {
			log.Errorf("run cmd: %s, err: %v", cmdStr, err)
			return err
		}
	}

	log.Info("start run pip3 install...")
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && pip3 install -r requirements.txt", gitDir))
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}
