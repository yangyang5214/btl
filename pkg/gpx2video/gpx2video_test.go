package gpx2video

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tkrajina/gpxgo/gpx"
	"testing"
)

func TestRouteVideo(t *testing.T) {
	gpxFilePath := "/Users/beer/Downloads/GOTOES_1421451299181451.gpx"

	//f, err := os.CreateTemp("", "gpx2video*")
	//if err != nil {
	//	panic(err)
	//}

	workDir := "/tmp/111"

	gpxData, err := gpx.ParseFile(gpxFilePath)
	if err != nil {
		panic(err)
	}
	route := NewRouteVideo(gpxData, log.DefaultLogger, workDir)
	err = route.Run()
	if err != nil {
		panic(err)
	}
}

func TestGenVideo(t *testing.T) {
	route := RouteVideo{
		log: log.NewHelper(log.DefaultLogger),
	}
	err := route.genVideo("/tmp/111", "/tmp/1.mp4")
	if err != nil {
		panic(err)
	}
}
