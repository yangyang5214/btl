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
