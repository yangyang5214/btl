package pkg

import (
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestShot(t *testing.T) {
	html := "file:///Users/beer/beer/btl/index.html"
	ss := NewScreenshot("result.png", html, log.DefaultLogger)
	err := ss.Run()
	if err != nil {
		t.Fatal(err)
	}
}
