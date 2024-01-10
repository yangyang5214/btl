package gpx_amap

import (
	"golang.org/x/image/colornames"
	"image/color"
	"testing"

	"github.com/yangyang5214/btl/pkg/utils"
)

func Test1(t *testing.T) {
	//files := utils.FindGpxFiles("/Users/beer/beer/rides/shang_hai")
	files := utils.FindGpxFiles("/tmp/test/")
	//files := utils.FindGpxFiles("/Users/beer/beer/gpx_export/garmin_export_out/287053469.gpx")

	//styles := []string{"whitesmoke", "grey", "dark", "light", "fresh", "blue", "darkblue", "macaron"}
	styles := []string{"8ee61a45840f14ac60f33a799fbd00d8"}
	for _, style := range styles {
		ga := NewGpxAmap(style)
		ga.SetFiles(files)
		ga.SetColors([]color.Color{
			colornames.White,
		})
		ga.SetStrokeWeight(10)
		ga.HideStartEndPoint()
		err := ga.Run()
		if err != nil {
			t.Fatal(err)
		}
	}
}
