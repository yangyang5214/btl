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

func TestFitStat(t *testing.T) {
	session, err := ParseFit("/Users/beer/beer/merge-fit/fits/1.fit", log.NewHelper(log.DefaultLogger))
	if err != nil {
		panic(err)
	}
	t.Log(session)
}

func TestFitImg(t *testing.T) {
	//session, err := ParseFit("/Users/beer/beer/merge-fit/fits/1.fit", log.NewHelper(log.DefaultLogger))

	gpxData, err := gpx.ParseFile("/Users/beer/Downloads/84712d7f25ad6d1fb2793f202ac82708.gpx")
	if err != nil {
		panic(err)
	}
	session, err := ParseGPX(gpxData)
	if err != nil {
		panic(err)
	}
	err = NewImgOverview(session, log.DefaultLogger).Run()
	if err != nil {
		panic(err)
	}
}

func TestImgFit(t *testing.T) {
}
