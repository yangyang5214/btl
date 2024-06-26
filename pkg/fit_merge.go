package pkg

import (
	"bytes"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yangyang5214/btl/pkg/utils"
	"os"
	"os/exec"
	"path"
	"strings"
)

type FitMerge struct {
	log *log.Helper
}

func NewFitMerge() *FitMerge {
	return &FitMerge{
		log: log.NewHelper(log.DefaultLogger),
	}
}

func (s *FitMerge) Run() error {
	fitFiles := utils.FindFitFiles(".")
	resultData, err := MergeFit(fitFiles, s.log)
	if err != nil {
		return err
	}
	return os.WriteFile("result.fit", resultData, 0755)
}

func MergeFit(fitFiles []string, logHelper *log.Helper) ([]byte, error) {
	jarEnv := MergeFitJarPath()
	workDir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Remove(workDir)
	}()
	resultFit := path.Join(workDir, uuid.New().String()+".fit")
	defer func() {
		_ = os.Remove(resultFit)
	}()
	gitMergeCmd := fmt.Sprintf("/usr/bin/java -jar %s -m %s %s", jarEnv, resultFit, strings.Join(fitFiles, " "))
	cmd := exec.Command("/bin/bash", "-c", gitMergeCmd)

	logHelper.Infof("run cmd: %s", gitMergeCmd)

	var logBuffer bytes.Buffer
	cmd.Stdout = &logBuffer
	cmd.Stderr = &logBuffer

	_ = cmd.Run()

	logHelper.Infof("run merge-fit: %s", logBuffer.String())
	_, err = os.Stat(resultFit)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return os.ReadFile(resultFit)
}

func MergeFitJarPath() string {
	homeDir, _ := os.UserHomeDir()
	return path.Join(homeDir, "merge-fit.jar")
}
