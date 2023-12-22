package pkg

import (
	"golang.org/x/image/colornames"
	"testing"
)

func TestColorToHex(t *testing.T) {
	c := colornames.Yellow
	r := ColorToHex(c)
	t.Log(r)
}

func TestFileExists(t *testing.T) {
	t.Log(FileExists("/tmp"))
	t.Log(FileExists("/tmp/sfsfsfsf"))
	t.Log(FileExists("/Users/beer/beer/gpx_export/garmin_export_out"))
}
