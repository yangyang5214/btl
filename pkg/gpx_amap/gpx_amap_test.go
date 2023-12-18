package gpx_amap

import (
	"github.com/yangyang5214/btl/pkg/utils"
	"golang.org/x/image/colornames"
	"image/color"
	"testing"
)

func Test1(t *testing.T) {
	files := utils.FindGpxFiles("/Users/beer/beer/rides/shang_hai")
	ga := NewGpxAmap(files)
	ga.SetColors([]color.Color{colornames.Red})
	//ga.SetMapStyle(Dark)
	//ga.SetMapStyle(Light)
	//ga.SetMapStyle(Whitesmoke)
	//ga.SetMapStyle(Grey)
	//ga.SetMapStyle(Fresh)
	//ga.SetMapStyle(Blue)
	//ga.SetMapStyle(Darkblue)
	//ga.SetMapStyle(Macaron)
	err := ga.Run()
	if err != nil {
		t.Fatal(err)
	}
}
