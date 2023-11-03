package pkg

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"sort"
	"sync"

	"github.com/yangyang5214/gou"

	"github.com/schollz/progressbar/v3"

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
		smCtx.AddObject(sm.NewPath(position, color.Transparent, 3))
	}
	return smCtx.DetermineZoomCenter()
}

func (g *GpxVideo) genStep(positions [][]s2.LatLng) int {
	var lines []int
	for _, position := range positions {
		lines = append(lines, len(position))
	}
	sort.Ints(lines)
	log.Infof("all lines points is %v", lines)
	// 取一个线路为 基准
	index := len(lines) / 5 * 3

	log.Infof("use lines index for step %d", index)
	step := lines[index] / (30 * 12) // ffmpeg -framerate 30 * 15 seconds
	if step < 1 {
		step = 1
	}
	return step
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

	var img image.Image
	var index int64
	var resultImg string //last image

	maxPoints := g.maxPoints(positions)
	bar := progressbar.New(maxPoints)
	step := g.genStep(positions)

	log.Infof("max points count: %d, set step %d", maxPoints, step)

	var wg sync.WaitGroup
	for i := step; i < maxPoints; i += step {
		index = index + 1
		wg.Add(1)
		imgPath := path.Join(tempDir, fmt.Sprintf("%d.png", index))
		resultImg = imgPath

		go func(imgPath string, i int, step int, lines [][]s2.LatLng) {
			defer func() {
				_ = bar.Add(step)
				wg.Done()
			}()
			smCtx := sm.NewContext()
			smCtx.SetSize(width, height)
			smCtx.SetTileProvider(g.titleProvider)
			smCtx.SetZoom(zoom)
			smCtx.SetCenter(center)

			for lineIndex, line := range lines {
				end := i
				if end > len(line) {
					end = len(line)
				}
				ps := line[0:end]
				curColor := utils.GetColor(lineIndex, g.cols)
				smCtx.AddObject(sm.NewPath(ps, curColor, g.weight))
			}

			img, err = smCtx.Render()
			if err != nil {
				log.Errorf("call sm.Render failed %+v", err)
			}
			err = gg.SavePNG(imgPath, img)
			if err != nil {
				log.Errorf("save png failed %+v", err)
			}
		}(imgPath, i, step, positions)
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
