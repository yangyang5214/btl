package gpx_export

import "testing"

func TestRun(t *testing.T) {
	username := ""
	password := ""
	g := NewGpxExport(GarminCN, username, password)
	err := g.Run()
	if err != nil {
		t.Fatal(err)
	}
}
