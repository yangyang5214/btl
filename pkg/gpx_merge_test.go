package pkg

import (
	"testing"
)

func Test_ParseTrkseg(t *testing.T) {
	g := NewGpxMerge("")
	r, err := g.parseTrkseg("/tmp/1.gpx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestMerge(t *testing.T) {
	g := NewGpxMerge("/tmp/1")
	err := g.Run()
	if err != nil {
		t.Fatal(err)
	}
}
