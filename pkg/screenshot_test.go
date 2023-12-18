package pkg

import "testing"

func TestShot(t *testing.T) {
	ss := NewScreenshot("result.png", "/Users/beer/beer/btl/pkg/gpx_amap/index.html")
	err := ss.Run()
	if err != nil {
		t.Fatal(err)
	}
}
