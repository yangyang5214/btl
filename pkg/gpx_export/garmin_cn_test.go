package gpx_export

import (
	"os"
	"path"
	"strings"
	"testing"
)

var (
	username string
	password string
)

func init() {
	homeDir, _ := os.UserHomeDir()
	r, err := os.ReadFile(path.Join(homeDir, ".garmin_cn"))
	if err != nil {
		panic(err)
	}
	content := string(r)
	lines := strings.Split(content, "\n")

	username = lines[0]
	password = lines[1]
}

func TestRun(t *testing.T) {
	g := NewGpxExport(GarminCN, username, password)
	err := g.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogin(t *testing.T) {
	g := NewGarminCn()
	t.Logf("username: <%v>, password: <%s>", username, password)
	g.Init("/Users/beer/.gpx_export", "/tmp", username, password)
	t.Log(g.Auth())
}
