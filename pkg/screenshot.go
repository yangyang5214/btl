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
	chromeCtx, cancel := chromedp.NewContext(ctxTimeout)
	defer cancel()

	var buf []byte
	if err = chromedp.Run(chromeCtx,
		chromedp.Navigate(s.urlStr),
		chromedp.EmulateViewport(1440, 900),
		chromedp.Sleep(s.waitSeconds*2),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		return errors.WithStack(err)
	}
	return os.WriteFile(s.imgPath, buf, 0644)
}
