package pkg

import (
	"image/color"
	"strconv"
	"testing"

	"golang.org/x/image/colornames"

	"github.com/fogleman/gg"
	sm "github.com/yangyang5214/go-staticmaps"
)

func Test(t *testing.T) {
	files := []string{
		//"/Users/beer/beer/rides/hai_nan.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_190230645.gpx",
		"/Users/beer/beer/rides/shang_hai/7198278844.gpx",
		"/Users/beer/beer/rides/shang_hai/6884490616.gpx",
	}
	names := []string{
		"none",
		//"carto-dark", //黑暗系
		//"arcgis-worldimagery", //地理影像
		//"carto-light", //白色系
		//"wikimedia", //纯线路
		//"stamen-terrain",      //森林？
	}
	for _, name := range names {
		gpx := NewGpxMap(files, "", name, []color.Color{colornames.White})
		gpx.SetBgColor(colornames.Yellow)
		gpx.SetWeight(5)
		err := gpx.Run("result.png")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestName(t *testing.T) {
	ctx := sm.NewContext()
	h := 1200
	w := float64(h) * 1.5

	ctx.SetSize(int(w), h)
	r, err := sm.ParsePathString("colors:red|weight:5|gpx:/Users/beer/Downloads/1.gpx")
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range r {
		ctx.AddObject(path)

		length := len(path.Positions)
		last := path.Positions[length-1]
		first := path.Positions[0]

		t.Log(last.Lat.String())
		t.Log(last.Distance(first))

		r, _ := strconv.ParseFloat(last.Lat.String(), 10)
		r1, _ := strconv.ParseFloat(first.Lat.String(), 10)

		t.Log(r1 - r)
	}

	img, err := ctx.Render()
	if err != nil {
		t.Fatal(err)
	}

	if err = gg.SavePNG("1.png", img); err != nil {
		t.Fatal(err)
	}
}

func TestColors(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	for i := 0; i < 20; i++ {
		t.Log(s[i%len(s)])
	}
}
