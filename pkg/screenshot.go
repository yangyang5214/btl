package pkg

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/pkg/errors"

	"github.com/chromedp/chromedp"
)

type Screenshot struct {
	imgPath     string
	urlStr      string
	waitSeconds time.Duration
	ctx         context.Context
	log         *log.Helper
}

func NewScreenshot(imgPath, urlStr string, logger log.Logger) *Screenshot {
	return &Screenshot{
		imgPath:     imgPath,
		urlStr:      urlStr,
		waitSeconds: 3 * time.Second,
		ctx:         context.Background(),
		log:         log.NewHelper(logger),
	}
}

func (s *Screenshot) SetWaitSeconds(wait int32) {
	s.waitSeconds = time.Duration(wait) * time.Second
}

// Run https://github.com/chromedp/chromedp/issues/941#issuecomment-961181348
func (s *Screenshot) Run() error {
	s.log.Infof("screenshot for %s, waitSeconds is %v", s.urlStr, s.waitSeconds)
	ctxTimeout, cancel := context.WithTimeout(s.ctx, time.Minute*3)
	defer cancel()

	dir := filepath.Dir(s.imgPath)
	_ = os.MkdirAll(dir, os.ModePerm)

	var err error
	s.log.Info("Creating Chrome context...")
	chromeCtx, cancel := chromedp.NewContext(ctxTimeout)
	defer cancel()

	var buf []byte
	s.log.Info("Starting chromedp.Run...")
	if err = chromedp.Run(chromeCtx,
		chromedp.Navigate(s.urlStr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			s.log.Info("Navigation complete")
			return nil
		}),
		chromedp.EmulateViewport(1440, 900),
		chromedp.ActionFunc(func(ctx context.Context) error {
			s.log.Info("Viewport emulation complete")
			return nil
		}),
		chromedp.Sleep(s.waitSeconds*2),
		chromedp.ActionFunc(func(ctx context.Context) error {
			s.log.Info("Sleep complete")
			return nil
		}),
		chromedp.CaptureScreenshot(&buf),
		chromedp.ActionFunc(func(ctx context.Context) error {
			s.log.Info("Screenshot capture complete")
			return nil
		}),
	); err != nil {
		s.log.Errorf("Error during chromedp.Run: %v", err)
		return errors.WithStack(err)
	}

	s.log.Info("Writing screenshot to file...")
	if err = os.WriteFile(s.imgPath, buf, 0644); err != nil {
		s.log.Errorf("Error writing file: %v", err)
		return errors.WithStack(err)
	}

	s.log.Info("Screenshot saved successfully")
	return nil
}
