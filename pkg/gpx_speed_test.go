package pkg

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func Test1(t *testing.T) {
	gs := NewGpxSpeed("./../example/demo.gpx", 100, log.DefaultLogger)
	err := gs.Run()
	if err != nil {
		panic(err)
	}
}
