package pkg

import (
	"context"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	time "time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Screenshot struct {
	imgPath     string
	htmlPath    string
	waitSeconds time.Duration
	ctx         context.Context
}

func NewScreenshot(imgPath, htmlPath string) *Screenshot {
	return &Screenshot{
		imgPath:     imgPath,
		htmlPath:    htmlPath,
		waitSeconds: 5 * time.Second,
		ctx:         context.Background(),
	}
}

func (s *Screenshot) SetWaitSeconds(wait int32) {
	s.waitSeconds = time.Duration(wait) * time.Second
}

// Run https://github.com/chromedp/chromedp/issues/941#issuecomment-961181348
func (s *Screenshot) Run() error {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*60)
	defer cancel()

	dir := filepath.Dir(s.imgPath)
	_ = os.MkdirAll(dir, os.ModePerm)

	html, err := os.ReadFile(s.htmlPath)
	if err != nil {
		return errors.WithStack(err)
	}
	chromeCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	if err = chromedp.Run(chromeCtx,
		// the navigation will trigger the "page.EventLoadEventFired" event too,
		// so we should add the listener after the navigation.
		chromedp.Navigate("about:blank"),
		chromedp.EmulateViewport(1440, 900),
		// set the page content and wait until the page is loaded (including its resources).
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return errors.WithStack(err)
			}

			if err := page.SetDocumentContent(frameTree.Frame.ID, string(html)).Do(ctx); err != nil {
				return errors.WithStack(err)
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(s.waitSeconds)
			cs := page.CaptureScreenshot()
			cs.Quality = 100
			buf, err := cs.Do(ctx)
			if err != nil {
				return errors.WithStack(err)
			}
			return os.WriteFile(s.imgPath, buf, 0644)
		}),
	); err != nil {
		return errors.WithStack(err)
	}
	return chromedp.Cancel(chromeCtx)
}
