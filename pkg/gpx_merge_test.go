package pkg

import (
	"os"
	"testing"

	"github.com/tkrajina/gpxgo/gpx"
)

func TestMerge(t *testing.T) {
	g := NewGpxMerge("/tmp/test")
	err := g.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestName(t *testing.T) {
	bytes, err := os.ReadFile("/Users/beer/Downloads/gpx-1/result.gpx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(bytes) / 1024 / 1024)
}

func TestSetGpxDatas(t *testing.T) {
	bytes, err := os.ReadFile("/tmp/111/13175318687.gpx")
	if err != nil {
		t.Fatal(err)
	}
	g := NewGpxMerge("")

	gpxData, err := gpx.ParseBytes(bytes)
	if err != nil {
		panic(err)
	}
	g.SetGpxDatas([]*gpx.GPX{gpxData})
	g.SetResultPath("/tmp/result.gpx")
	err = g.Run()
	if err != nil {
		t.Fatal(err)
	}
}
