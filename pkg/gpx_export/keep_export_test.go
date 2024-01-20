package gpx_export

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestKeepLogin(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	r, err := os.ReadFile(path.Join(homeDir, ".keep"))

	datas := string(r)
	lines := strings.Split(datas, "\n")
	mobile := lines[0]
	pwd := lines[1]
	success, err := login(mobile, pwd)
	if err != nil {
		panic(err)
	}
	t.Log(success)
}
