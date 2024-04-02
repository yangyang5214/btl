package pkg

import (
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestShot(t *testing.T) {
	ss := NewScreenshot("result.png", "/Users/beer/beer/btl/pkg/gpx_amap/index.html", log.DefaultLogger)
	err := ss.Run()
	if err != nil {
		t.Fatal(err)
	}
}
