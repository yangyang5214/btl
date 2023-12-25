package pkg

import (
	"os"
	"testing"
)

func TestMerge(t *testing.T) {
	g := NewGpxMerge("/tmp/2")
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
	err = g.SetGpxDatas([][]byte{bytes})
	if err != nil {
		t.Fatal(err)
	}
	g.SetResultPath("/tmp/result.gpx")
	err = g.Run()
	if err != nil {
		t.Fatal(err)
	}
}
