package merge_img

import (
	"bufio"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	"image"
	_ "image/gif"  // Import gif support
	_ "image/jpeg" // Import jpeg support
	"image/png"
	_ "image/png" // Import png support
	"os"
)

type MergeImg struct {
	fgImgPath string
	bgImgPath string

	fgImg image.Image
	bgImg image.Image
	log   *log.Helper
}

func NewMergeImg(fgImgPath, bgImgPath string, logger log.Logger) *MergeImg {
	return &MergeImg{
		fgImgPath: fgImgPath,
		bgImgPath: bgImgPath,
		log:       log.NewHelper(log.With(logger, "btl-model", "merge_img")),
	}
}

func (s *MergeImg) init() error {
	// Open and decode the foreground image
	fgFile, err := os.Open(s.fgImgPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open foreground image: %s", s.fgImgPath)
	}
	defer fgFile.Close()
	s.fgImg, _, err = image.Decode(fgFile)
	if err != nil {
		return errors.Wrapf(err, "failed to decode foreground image: %s", s.fgImgPath)
	}

	// Open and decode the background image
	bgFile, err := os.Open(s.bgImgPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open background image: %s", s.bgImgPath)
	}
	defer bgFile.Close()
	s.bgImg, _, err = image.Decode(bufio.NewReader(bgFile))
	if err != nil {
		return errors.Wrapf(err, "failed to decode background image: %s", s.bgImgPath)
	}
	return nil
}

func (s *MergeImg) process() (*image.RGBA, error) {
	err := s.init()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Get the bounds of the background image
	bgBounds := s.bgImg.Bounds()
	newWidth := int(float64(bgBounds.Dx()) / 5 * 3)

	ratio := float64(newWidth) / float64(s.fgImg.Bounds().Dx())
	s.log.Infof("ratio is %f", ratio)

	newHeight := int(ratio * float64(s.fgImg.Bounds().Dy()))

	// Create a new image with the same size as the resized foreground image
	resizedFgImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Resize the foreground image
	draw.CatmullRom.Scale(resizedFgImg, resizedFgImg.Bounds(), s.fgImg, s.fgImg.Bounds(), draw.Over, nil)

	// Create a new image with the same size as the background
	mergedImg := image.NewRGBA(bgBounds)

	// Draw the background image onto the new image
	draw.Draw(mergedImg, bgBounds, s.bgImg, image.Point{}, draw.Src)

	// Calculate the position for the resized foreground image (right upper corner)
	offset := image.Pt(bgBounds.Dx()-newWidth, 0)
	fgBounds := image.Rect(offset.X, offset.Y, offset.X+newWidth, offset.Y+newHeight)

	// Draw the resized foreground image onto the new image
	draw.Draw(mergedImg, fgBounds, resizedFgImg, image.Point{}, draw.Over)

	return mergedImg, nil
}

func (s *MergeImg) Run(resultPath string) error {
	mergedImg, err := s.process()
	if err != nil {
		return errors.WithStack(err)
	}

	outFile, err := os.Create(resultPath)
	if err != nil {
		return errors.Wrap(err, "failed to create result image file")
	}
	defer outFile.Close()

	err = png.Encode(outFile, mergedImg)
	if err != nil {
		return errors.Wrap(err, "failed to encode result image as PNG")
	}

	return nil
}
