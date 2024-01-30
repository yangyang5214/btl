package utils

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/tkrajina/gpxgo/gpx"
)

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

func TestIndent(t *testing.T) {
	f, err := gpx.ParseFile("1.gpx")
	if err != nil {
		panic(err)
	}
	bytes, err := gpx.ToXml(f, gpx.ToXmlParams{
		Indent: true,
	})
	if err != nil {
		panic(err)
	}

	resultFile, err := os.Create("result.gpx")
	defer resultFile.Close()
	if err != nil {
		panic(err)
	}
	_, err = resultFile.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func TestPoint(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		f := 110.0
		t.Logf("result is %v", strings.TrimRight(fmt.Sprintf("%.10f", f), "0."))  //  11
		t.Logf("result is %v", strings.TrimSuffix(fmt.Sprintf("%.10f", f), "0.")) // 110.0000000000
	})

	t.Run("test2", func(t *testing.T) {
		f := 0.000056
		t.Logf("result is %v", strings.TrimRight(fmt.Sprintf("%.10f", f), "0."))  //0.000056
		t.Logf("result is %v", strings.TrimSuffix(fmt.Sprintf("%.10f", f), "0.")) //0.000056
	})
}

func TestNewGpx(t *testing.T) {
	gpxData := &gpx.GPX{}
	bytes, err := gpx.ToXml(gpxData, gpx.ToXmlParams{
		Indent: true,
	})
	if err != nil {
		panic(err)
	}

	resultFile, err := os.Create("new.gpx")
	defer resultFile.Close()
	if err != nil {
		panic(err)
	}
	_, err = resultFile.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func Test1111(t *testing.T) {
	f, _ := gpx.ParseFile("/Users/beer/Downloads/HaiNan.gpx")
	speedsDistances := f.MovingData().SpeedsDistances
	t.Log(speedsDistances[0])
	t.Log(f.MovingData().MaxSpeed * 60 * 60 / 1000.0)
	t.Log(f.Tracks[0].MovingData().MaxSpeed * 60 * 60 / 1000.0)
	t.Log(f.Tracks[0].Segments[0].MovingData().MaxSpeed * 60 * 60 / 1000.0)
}
