package pkg

import (
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	files := []string{
		"/Users/beer/beer/rides/shang_hai/activity_188416039.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_190230645.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_192638282.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_195179761.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_198260918.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_198428371.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_198704961.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_201117043.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_201386247.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_221242069.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_226372078.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_226427688.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_226460680.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_226509692.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_226542964.gpx",
		"/Users/beer/beer/rides/shang_hai/activity_226588901.gpx",
		"/Users/beer/beer/rides/shang_hai/6789896328.gpx",
		"/Users/beer/beer/rides/shang_hai/7198278844.gpx",
		"/Users/beer/beer/rides/shang_hai/6884490616.gpx",
		"/Users/beer/beer/rides/shang_hai/6830932310.gpx",
		"/Users/beer/beer/rides/shang_hai/6833367971.gpx",
	}
	names := []string{
		//"carto-dark",          //黑暗系
		//"arcgis-worldimagery", //地理影像
		"carto-light", //白色系
		//"wikimedia", //纯线路
		//"stamen-terrain",      //森林？
	}
	for _, name := range names {
		gpx := NewGpxMap(files, "", name, nil)
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
	r, err := sm.ParsePathString("color:red|weight:5|gpx:/Users/beer/Downloads/1.gpx")
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
