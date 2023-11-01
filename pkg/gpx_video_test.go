package pkg

import "testing"

func TestVideo(t *testing.T) {
	files := []string{
		"/Users/beer/beer/rides/shang_hai/activity_188416039.gpx",
	}
	v := NewGpxVideo(files)
	err := v.Run()
	if err != nil {
		t.Fatal(err)
	}
}
