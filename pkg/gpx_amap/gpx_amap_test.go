package gpx_amap

import (
	"fmt"
	"image/color"
	"testing"

	"golang.org/x/image/colornames"

	"github.com/yangyang5214/btl/pkg/utils"
)

func Test1(t *testing.T) {
	files := utils.FindGpxFiles("/Users/beer/beer/rides/shang_hai")
	//files := utils.FindGpxFiles("/Users/beer/beer/gpx_export/garmin_export_out/287053469.gpx")

	//styles := []string{"whitesmoke", "grey", "dark", "light", "fresh", "blue", "darkblue", "macaron"}
	styles := []string{"darkblue"}
	for _, style := range styles {
		ga := NewGpxAmap(style)
		ga.SetFiles(files)
		ga.SetColors([]color.Color{
			colornames.White,
		})
		ga.SetStrokeWeight(10)
		ga.HideStartEndPoint()
		ga.SetImgPath(fmt.Sprintf("/tmp/result/%s.png", style))
		err := ga.Run()
		if err != nil {
			t.Fatal(err)
		}
	}
}
