package pkg

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
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
	col           color.Color
	titleProvider *sm.TileProvider

	pointCount int

	pwd    string
	weight float64
}

var (
	outMp4 = "out.mp4"
)

func NewGpxVideo(files []string) (*GpxVideo, error) {
	if !gou.CmdExists("ffmpeg") {
		return nil, errors.New("ffmpeg not found")
	}

	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}
	col, _ := sm.ParseColorString("green")
	titleProvider := sm.NewTileProviderCartoDark()
	titleProvider.Attribution = ""

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &GpxVideo{
		files:         files,
		col:           col,
		titleProvider: titleProvider,
		pwd:           pwd,
		weight:        2,
	}, nil
}

func (g *GpxVideo) genCenter(positions [][]s2.LatLng) (int, s2.LatLng, error) {
	smCtx := sm.NewContext()
	for _, position := range positions {
		smCtx.AddObject(sm.NewPath(position, g.col, 2))
	}
	return smCtx.DetermineZoomCenter()
}

func (g *GpxVideo) genStep() int {
	step := g.pointCount / (30 * 10) // ffmpeg -framerate 30 * 15 seconds
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
	g.pointCount = utils.CountPoints(positions)
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

	bar := progressbar.New(g.pointCount)
	step := g.genStep()

	log.Infof("points count: %d, set step %d", g.pointCount, step)

	var wg sync.WaitGroup
	for _, position := range positions {
		for i := step; i < len(position); i += step {
			index = index + 1
			wg.Add(1)
			imgPath := path.Join(tempDir, fmt.Sprintf("%d.png", index))
			resultImg = imgPath
			go func(imgPath string, i int, step int) {
				defer func() {
					_ = bar.Add(step)
					wg.Done()
				}()
				smCtx := sm.NewContext()
				smCtx.SetSize(width, height)
				smCtx.SetTileProvider(g.titleProvider)
				smCtx.SetZoom(zoom)
				smCtx.SetCenter(center)

				smCtx.AddObject(sm.NewPath(position[0:i], g.col, g.weight))

				img, err = smCtx.Render()
				if err != nil {
					log.Errorf("call sm.Render failed %+v", err)
				}
				err = gg.SavePNG(imgPath, img)
				if err != nil {
					log.Errorf("save png failed %+v", err)
				}
			}(imgPath, i, step)
		}
	}
	wg.Wait()

	_ = bar.Finish() //set finished

	log.Info("\n")
	log.Infof("Satrt merge to  video ..")
	err = fileutil.CopyFile(resultImg, path.Join(g.pwd, "result.png"))
	if err != nil {
		return errors.WithStack(err)
	}
	cmd := "ffmpeg -framerate 30 -i %d.png -c:v libx264  -y " + outMp4
	err, _ = gou.RunCmd(fmt.Sprintf("cd %s", tempDir) + " && " + cmd)
	if err != nil {
		log.Errorf("run ffmpeg error %v", err)
		return errors.WithStack(err)
	}

	err = fileutil.CopyFile(path.Join(tempDir, outMp4), g.pwd)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
