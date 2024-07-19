package pkg

import (
	"github.com/fogleman/gg"
	"github.com/go-kratos/kratos/v2/log"
	"image"
	"image/color"
	"os"
	"strings"
)

type ImgStat struct {
	imgPath string
	log     *log.Helper
	result  string
}

type StatInfo struct {
	Label string
	Value string
}

func (s *StatInfo) String() string {
	return s.Label + ":" + s.Value
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func NewImgStat(imgPath string, logger log.Logger) *ImgStat {
	return &ImgStat{
		imgPath: imgPath,
		log:     log.NewHelper(logger),
		result:  "result.png",
	}
}

func (s *ImgStat) SetResult(result string) {
	s.result = result
}
func (s *ImgStat) Run(infos []*StatInfo) error {
	s.log.Infof("use stat info %v", infos)

	img, err := loadImage(s.imgPath)
	if err != nil {
		return err
	}
	dc := gg.NewContextForImage(img)

	height := float64(dc.Height())
	width := float64(dc.Width())

	fontSize := height / 50
	s.log.Infof("fontSize use %v", fontSize)
	lineSpacing := fontSize + fontSize*1.5
	x := fontSize * 2
	y := fontSize * 3

	for _, info := range infos {
		drawStr(dc, fontSize, info.Label, x, y)
		y += lineSpacing
		drawStr(dc, fontSize*2, strings.ToUpper(info.Value), x, y)
		y += lineSpacing
	}
	dc.Stroke()

	//add gpxt
	//textX := width - (width/fontSize)*1.5
	//textY := height - height/fontSize
	//dc.SetColor(colornames.Red)
	//_ = LoadFontFace(dc, fontSize*2)
	//dc.DrawStringAnchored("gpxt", textX, textY, 0.5, 0.5)

	return dc.SavePNG(s.result)
}

func drawStr(dc *gg.Context, fontSize float64, label string, x, y float64) {
	_ = LoadFontFace(dc, fontSize)
	dc.SetColor(color.White)
	dc.DrawStringAnchored(label, x, y, 0, 0)
}
