package pkg

import (
	"testing"
	"time"

	"github.com/schollz/progressbar/v3"
)

func TestVideo(t *testing.T) {
	files := []string{
		"/Users/beer/beer/rides/shang_hai/activity_188416039.gpx",
	}
	v, err := NewGpxVideo(files, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = v.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestProcess(t *testing.T) {
	bar := progressbar.Default(100)
	for i := 0; i < 100; i++ {
		_ = bar.Add(1)
		time.Sleep(100 * time.Millisecond)
	}
}
