package pkg

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"path"
	"testing"
)

// 批量 fit 转 http
func TestFit2Gpx(t *testing.T) {
	targetDir := "/Users/beer/Downloads/5月2日-5月26日318川藏骑行数据"
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		panic(err)
	}
	for index, entry := range entries {
		t.Logf(entry.Name())

		p := path.Join(targetDir, entry.Name())
		fp := NewFit2Gpx(p, log.DefaultLogger)
		fp.SetResultPath(fmt.Sprintf("/tmp/gpxs/%d.gpx", index))
		err = fp.Run()
		if err != nil {
			panic(err)
		}
	}
}

func TestFit2Csv(t *testing.T) {
	fg := NewFit2Gpx("/Users/beer/Downloads/240525144003.fit", log.DefaultLogger)
	err := fg.fit2Csv("/tmp/1.csv")
	if err != nil {
		panic(err)
	}
}

func TestParserCsv(t *testing.T) {
	fg := NewFit2Gpx("/Users/beer/Downloads/240525144003.fit", log.DefaultLogger)
	session, err := fg.parserCsv("/tmp/1.csv")
	if err != nil {
		panic(err)
	}
	t.Log(session.points[0])
}

func TestNewFit2Gpx(t *testing.T) {
	fp := NewFit2Gpx("/Users/beer/Downloads/Night_Run.fit", log.DefaultLogger)
	err := fp.Run()
	if err != nil {
		panic(err)
	}
}
