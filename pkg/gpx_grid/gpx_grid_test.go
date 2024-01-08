package gpx_grid_test

import (
	"testing"

	"github.com/yangyang5214/btl/pkg/gpx_grid"
	"github.com/yangyang5214/btl/pkg/utils"
)

func TestName(t *testing.T) {
	files := utils.FindGpxFiles("/Users/beer/beer/rides/shang_hai")
	g := gpx_grid.NewGpxGrid()
	g.SetFiles(files)
	err := g.Run()
	if err != nil {
		panic(err)
	}
}
