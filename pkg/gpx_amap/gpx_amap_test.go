package gpx_amap

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/yangyang5214/btl/pkg/utils"
	"golang.org/x/image/colornames"
)

func Test1(t *testing.T) {
	files := utils.FindGpxFiles("/Users/beer/beer/rides/shang_hai")

	styles := []string{"whitesmoke", "grey", "dark", "light", "fresh", "blue", "darkblue", "macaron"}
	for _, style := range styles {
		ga := NewGpxAmap(style)
		ga.SetFiles(files)
		ga.SetStep(20)
		ga.SetColors([]color.Color{
			colornames.Red,
		})
		ga.SetImgPath(fmt.Sprintf("/tmp/result/%s.png", style))
		err := ga.Run()
		if err != nil {
			t.Fatal(err)
		}
	}
}
