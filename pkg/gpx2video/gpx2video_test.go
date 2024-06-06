package gpx2video

import "testing"

func TestName(t *testing.T) {
	gpxFilePath := "/Users/beer/Downloads/activity_365437876.gpx"

	points, err := parseGPX(gpxFilePath)
	if err != nil {
		panic(err)
	}

	err = plotGPX(points, "result.png")
	if err != nil {
		panic(err)
	}
}
