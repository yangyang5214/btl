package pkg

import (
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	files := []string{
		//"/Users/beer/Downloads/1.gpx",
		//"/Users/beer/Downloads/activity_273093158.gpx",
		//"/Users/beer/Downloads/activity_190230645.gpx",
		//"/Users/beer/Downloads/activity_198260918.gpx",
		"/Users/beer/Downloads/HaiNan.gpx",
	}
	gpx := NewGpxMap(files, "beer")
	err := gpx.Run("result.png")
	if err != nil {
		t.Fatal(err)
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
