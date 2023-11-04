package pkg

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/yangyang5214/gou"
	"image/color"
	"os"
	"path"
	"sync"

	fileutil "github.com/yangyang5214/gou/file"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg/utils"
)

type GpxVideo struct {
	files         []string
	cols          []color.Color
	titleProvider *sm.TileProvider

	pwd    string
	weight float64
}

var (
	outMp4 = "out.mp4"
)

func NewGpxVideo(files []string, cols []color.Color) (*GpxVideo, error) {
	if !gou.CmdExists("ffmpeg") {
		return nil, errors.New("ffmpeg not found")
	}

	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}
	titleProvider := sm.NewTileProviderCartoDark()
	titleProvider.Attribution = ""

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &GpxVideo{
		files:         files,
		cols:          cols,
		titleProvider: titleProvider,
		pwd:           pwd,
		weight:        2,
	}, nil
}

func (g *GpxVideo) genCenter(positions [][]s2.LatLng) (int, s2.LatLng, error) {
	smCtx := sm.NewContext()
	for _, position := range positions {
		smCtx.AddObject(sm.NewPath(position, color.Transparent, 1))
	}
	return smCtx.DetermineZoomCenter()
}

func (g *GpxVideo) genStep(positions [][]s2.LatLng) map[int]int {
	r := make(map[int]int)
	for index := range positions {
		//step := len(position)/(30*15) + 1
		//if step > 10 {
		//	step = 10
		//}
		r[index] = 5
	}
	return r
}

func (g *GpxVideo) initSmCtx(width, height, zoom int, center s2.LatLng) *sm.Context {
	smCtx := sm.NewContext()
	smCtx.SetSize(width, height)
	smCtx.SetTileProvider(g.titleProvider)
	smCtx.SetZoom(zoom)
	smCtx.SetCenter(center)
	return smCtx
}

func (g *GpxVideo) stepImage(bar *progressbar.ProgressBar, wg *sync.WaitGroup, imgPath string, positions [][]s2.LatLng, smCtx *sm.Context, index int, lineStep map[int]int) {
	defer func() {
		wg.Done()
		_ = bar.Add(1)
	}()

	for lineIndex, line := range positions {
		end := index * lineStep[lineIndex]
		if end > len(line) {
			end = len(line)
		}
		//log.Infof("size is %d,index is %d, end is %d", len(line), index, end)
		ps := line[0:end]
		curColor := utils.GetColor(lineIndex, g.cols)
		smCtx.AddObject(sm.NewPath(ps, curColor, g.weight))
	}

	img, err := smCtx.Render()
	if err != nil {
		log.Errorf("call sm.Render failed %+v", err)
	}
	err = gg.SavePNG(imgPath, img)
	if err != nil {
		log.Errorf("save png failed %+v", err)
	}

}

func (g *GpxVideo) Run() error {
	gpxDatas, err := utils.ParseGpxData(g.files)
	if err != nil {
		return err
	}
	positions, err := utils.ParsePositions(gpxDatas)
	if err != nil {
		return errors.WithStack(err)
	}
	zoom, center, err := g.genCenter(positions)
	if err != nil {
		return errors.WithStack(err)
	}

	width, height := utils.GenWidthHeight(positions)
	log.Infof("use height=%d, width=%d", height, width)

	tempDir, err := os.MkdirTemp("", uuid.New().String())
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	var index int
	var resultImg string //last image

	lineStep := g.genStep(positions)
	loops := g.loopCount(lineStep, positions)

	log.Infof("lineStep is %+v", lineStep)

	bar := progressbar.New(loops)

	wg := &sync.WaitGroup{}
	for i := 0; i < loops; i++ {
		index = index + 1
		wg.Add(1)
		imgPath := path.Join(tempDir, fmt.Sprintf("%d.png", index))
		resultImg = imgPath

		smCtx := g.initSmCtx(width, height, zoom, center)
		go g.stepImage(bar, wg, imgPath, positions, smCtx, index, lineStep)
	}
	wg.Wait()

	_ = bar.Finish() //set finished
	log.Info("\n")
	err = fileutil.CopyFile(resultImg, path.Join(g.pwd, "result.png"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = g.mergeVideo(tempDir)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (g *GpxVideo) maxPoints(points [][]s2.LatLng) int {
	var r int
	for _, point := range points {
		if len(point) > r {
			r = len(point)
		}
	}
	return r
}

func (g *GpxVideo) loopCount(points map[int]int, positions [][]s2.LatLng) int {
	var r int
	for index, step := range points {
		c := len(positions[index]) / step
		if c > r {
			r = c
		}
	}
	return r
}

func (g *GpxVideo) mergeVideo(workDir string) (err error) {
	log.Infof("Satrt merge to  video ..")

	cmd := "ffmpeg -framerate 30 -i %d.png -c:v libx264  -y " + outMp4
	err, _ = gou.RunCmd(fmt.Sprintf("cd %s", workDir) + " && " + cmd)
	if err != nil {
		log.Errorf("run ffmpeg error %v", err)
		return errors.WithStack(err)
	}

	err = fileutil.CopyFile(path.Join(workDir, outMp4), g.pwd)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
