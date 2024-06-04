package gpx_export

import (
	"fmt"
	"github.com/pkg/errors"
	"os/exec"

	"github.com/go-kratos/kratos/v2/log"
)

type GarminUpload struct {
	user     string
	pwd      string
	filepath string
	isCn     bool
	gitDir   string
	log      *log.Helper
}

func NewGarminUpload(user, pwd, filepath string, isCn bool, logger log.Logger) *GarminUpload {
	return &GarminUpload{
		user:     user,
		pwd:      pwd,
		filepath: filepath,
		isCn:     isCn,
		gitDir:   DefaultGitDir(),
		log:      log.NewHelper(logger),
	}
}

func (s *GarminUpload) Run() error {
	if s.user == "" {
		return fmt.Errorf("user is empty")
	}
	if s.pwd == "" {
		return fmt.Errorf("pwd is empty")
	}
	if s.filepath == "" {
		return fmt.Errorf("filepath is empty")
	}
	cmdStr := fmt.Sprintf("python3 %s/upload/garmin_upload.py -u '%s' -p '%s' -f %s",
		s.gitDir, s.user, s.pwd, s.filepath)

	if s.isCn {
		cmdStr = cmdStr + " --is-cn"
	}
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	s.log.Infof("start run cmd: %s", cmdStr)

	out, err := cmd.Output()
	if err != nil {
		return errors.New("导入失败")
	}
	s.log.Infof("run cmd out: %s", out)
	return nil
}
