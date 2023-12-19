package utils

import "testing"

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
