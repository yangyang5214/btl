package utils

import (
	"github.com/tkrajina/gpxgo/gpx"
	"testing"
)

func TestGetDate(t *testing.T) {
	r, err := ParseGpxData([]string{
		"/tmp/2/20191228_上午骑车.gpx",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(GetStartTime(r[0]))
}

func TestGpxtract(t *testing.T) {
	r, err := ParseGpxData([]string{
		"/tmp/2/20191228_上午骑车.gpx",
	})
	if err != nil {
		t.Fatal(err)
	}
	gpx := r[0]
	t.Log(gpx)
}

func TestName(t *testing.T) {
	f, _ := gpx.ParseFile("/tmp/test/result.gpx")
	for _, track := range f.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				if point.Longitude < 100 {
					t.Logf("%v,%v", point.Longitude, point.Latitude)
					t.Log(point)
				}
			}
		}
	}
}

func TestDistance(t *testing.T) {
	p1 := gpx.Point{
		Latitude:  18.4400882721,
		Longitude: 11,
	}
	p2 := gpx.Point{
		Latitude:  18.4400882721,
		Longitude: 110.3556594849,
	}

	distance := gpx.Distance2D(p1.Latitude, p1.Longitude, p2.Latitude, p2.Longitude, false)
	t.Log(distance)
}
