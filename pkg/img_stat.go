package pkg

import (
	"github.com/fogleman/gg"
	"github.com/go-kratos/kratos/v2/log"
	"image"
	"image/color"
	"os"
)

type ImgStat struct {
	imgPath string
	log     *log.Helper
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
	}
}

func (s *ImgStat) Run(infos []*StatInfo) error {
	s.log.Infof("use stat info %v", infos)

	img, err := loadImage(s.imgPath)
	if err != nil {
		return err
	}
	dc := gg.NewContextForImage(img)

	const fontSize = 20
	const lineSpacing = fontSize + fontSize
	x := 50.0
	y := 80.0

	for _, info := range infos {
		drawStr(dc, fontSize, info.Label, x, y)
		y += lineSpacing
		drawStr(dc, fontSize+10, info.Value, x, y)
		y += lineSpacing
	}
	dc.Stroke()
	return dc.SavePNG("result.png")
}

func drawStr(dc *gg.Context, fontSize float64, label string, x, y float64) {
	_ = LoadFontFace(dc, fontSize)
	dc.SetColor(color.White)
	dc.DrawStringAnchored(label, x, y, 0, 0)
}
