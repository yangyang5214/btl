package pkg

import (
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestChromePool_ScreenShot(t *testing.T) {
	chrome, cancel := NewChromePool(log.DefaultLogger)
	defer cancel()
	err := chrome.ScreenShot("https://www.baidu.com/", "/tmp/1.png")
	if err != nil {
		panic(err)
	}
}
