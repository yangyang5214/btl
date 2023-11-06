package pkg

import (
	"fmt"
	"image/color"
	"os"
	"path"
	"sort"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/yangyang5214/gou"

	fileutil "github.com/yangyang5214/gou/file"

	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg/utils"
	sm "github.com/yangyang5214/go-staticmaps"
)

type GpxVideo struct {
	files         []string
	cols          []color.Color
	titleProvider *sm.TileProvider

	pwd    string
	weight float64

	zoom   int
	center s2.LatLng
	width  int
	height int
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
		weight:        1,
	}, nil
}

func (g *GpxVideo) genCenter(positions [][]s2.LatLng) (int, s2.LatLng, error) {
	smCtx := sm.NewContext()
	for _, position := range positions {
		smCtx.AddObject(sm.NewPath(position, color.Transparent, g.weight))
	}
	return smCtx.DetermineZoomCenter()
}

func (g *GpxVideo) genStep(positions [][]s2.LatLng) map[int]int {
	var lines []int
	for _, position := range positions {
		lines = append(lines, len(position))
	}
	sort.Ints(lines)

	step := make(map[int]int)
	for index := range lines {
		step[index] = 5 //todo
	}
	return step
}

func (g *GpxVideo) initSmCtx() *sm.Context {
	smCtx := sm.NewContext()
	smCtx.SetSize(g.width, g.height)
	smCtx.SetTileProvider(g.titleProvider)
	smCtx.SetZoom(g.zoom)
	smCtx.SetCenter(g.center)
	return smCtx
}

func (g *GpxVideo) stepImage(bar *progressbar.ProgressBar, wg *sync.WaitGroup, imgPath string, positions [][]s2.LatLng, index int, lineStep map[int]int) {
	defer func() {
		_ = bar.Add(1)
		wg.Done()
	}()

	smCtx := g.initSmCtx()

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
	g.zoom = zoom
	g.center = center
	log.Infof("use zoom %d", zoom)

	width, height := utils.GenWidthHeight(positions)
	log.Infof("use height=%d, width=%d", height, width)
	g.width = width
	g.height = height

	tempDir, err := os.MkdirTemp("", uuid.New().String())
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	lineStep := g.genStep(positions)
	loops := g.loopCount(lineStep, positions) //最大 loops

	log.Infof("lineStep is %+v, loops is %d", lineStep, loops)

	bar := progressbar.New(loops)

	var index int
	wg := &sync.WaitGroup{}
	for i := 0; i < loops; i++ {
		index = index + 1
		wg.Add(1)
		imgPath := path.Join(tempDir, fmt.Sprintf("%d.png", index))
		go g.stepImage(bar, wg, imgPath, positions, index, lineStep)
	}
	wg.Wait()

	_ = bar.Finish() //set finished
	log.Info("\n")

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
