package gpx_export

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type KeepExporter struct {
	gitDir string

	exportDir string

	username string
	password string
	log      *log.Helper
}

const (
	LoginApi = "https://api.gotokeep.com/v1.1/users/login"
)

func login(mobile string, password string) (bool, error) {
	session := http.DefaultClient
	headers := http.Header{}
	headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	headers.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	data := url.Values{}
	data.Set("mobile", mobile)
	data.Set("password", password)

	resp, err := session.PostForm(LoginApi, data)
	if err != nil {
		return false, err
	}
	return resp.StatusCode == http.StatusOK, nil
}

func NewKeep(logger log.Logger) *KeepExporter {
	return &KeepExporter{
		log: log.NewHelper(log.With(logger, "app", "keep")),
	}
}

func (k *KeepExporter) Init(gitDir, exportDir string, username, password string) {
	k.exportDir = exportDir
	k.gitDir = gitDir
	k.username = username
	k.password = password
}

func (k *KeepExporter) Run() error {
	cmdStr := fmt.Sprintf("python3 %s/keep_export.py %s %s --out %s", k.gitDir, k.username, k.password, k.exportDir)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	k.log.Infof("start run cmd: %s", cmdStr)

	runtimeLog, err := cmd.Output()
	if err != nil {
		k.log.Infof("run cmd: %s", err)
		return err
	}
	result := string(runtimeLog)
	result = strings.Trim(result, "")
	k.log.Infof("run keep result:\n %s", result)

	cyclingDir := path.Join(k.exportDir, "cycling")
	_ = os.MkdirAll(cyclingDir, 0755)
	cmdStr = fmt.Sprintf("python3 %s/keep_export_cycling.py %s %s --out %s", k.gitDir, k.username, k.password, cyclingDir)
	cmd = exec.Command("/bin/bash", "-c", cmdStr)
	k.log.Infof("start run cmd: %s", cmdStr)

	runtimeLog, err = cmd.Output()
	if err != nil {
		k.log.Infof("run cmd: %s", err)
		return err
	}
	result = string(runtimeLog)
	result = strings.Trim(result, "")
	k.log.Infof("run keep result:\n %s", result)
	return nil
}

func (k *KeepExporter) Auth() bool {
	if len(k.username) != 11 {
		return false // 手机号
	}
	success, err := login(k.username, k.password)
	if err != nil {
		k.log.Errorf("login error: %v", err)
		return false
	}
	return success
}
