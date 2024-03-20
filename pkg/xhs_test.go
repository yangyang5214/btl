package pkg

import (
	"os"
	"testing"
)

func TestXpathUrl(t *testing.T) {
	bytes, err := os.ReadFile("/Users/beer/beer/btl/index.html")
	if err != nil {
		panic(err)
	}
	r, err := xpathUrl(string(bytes))
	if err != nil {
		panic(err)
	}
	t.Log(r.ImageUrls)
}
