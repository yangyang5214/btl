package pkg

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Screenshot struct {
	imgPath  string
	htmlPath string
}

func NewScreenshot(imgPath, htmlPath string) *Screenshot {
	return &Screenshot{
		imgPath:  imgPath,
		htmlPath: htmlPath,
	}
}

// Run https://github.com/chromedp/chromedp/issues/941#issuecomment-961181348
func (s *Screenshot) Run() error {
	dir := filepath.Dir(s.imgPath)
	_ = os.MkdirAll(dir, os.ModePerm)

	html, err := os.ReadFile(s.htmlPath)
	if err != nil {
		return err
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	if err = chromedp.Run(ctx,
		// the navigation will trigger the "page.EventLoadEventFired" event too,
		// so we should add the listener after the navigation.
		chromedp.Navigate("about:blank"),
		chromedp.EmulateViewport(1440, 900),
		// set the page content and wait until the page is loaded (including its resources).
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			if err := page.SetDocumentContent(frameTree.Frame.ID, string(html)).Do(ctx); err != nil {
				return err
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(3 * time.Second)
			buf, err := page.CaptureScreenshot().Do(ctx)
			if err != nil {
				return err
			}
			return os.WriteFile(s.imgPath, buf, 0644)
		}),
	); err != nil {
		return err
	}
	return nil
}
