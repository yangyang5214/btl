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
