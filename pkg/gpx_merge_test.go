package pkg

import (
	"testing"
)

func TestGarminGpx_getDate(t *testing.T) {
	g := NewGpxMerge("")
	got, err := g.getDate("/tmp/1.gpx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(got)
}

func Test_ParseTrkseg(t *testing.T) {
	g := NewGpxMerge("")
	r := g.parseTrkseg("/tmp/1.gpx")
	t.Log(r)
}
