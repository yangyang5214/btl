package pkg

import (
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestGenKml(t *testing.T) {
	opts := []Option{
		WithName("test"),
		WithResultFile("result.kml"),
	}
	k := NewGpx2Kml("/Users/beer/beer/btl/example/demo.gpx", log.DefaultLogger, opts...)
	err := k.Run()
	if err != nil {
		panic(err)
	}
}
