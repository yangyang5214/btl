package pkg

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
	Language  = "zh-CN,zh;q=0.9"
)

type ChromePool struct {
	browser *rod.Browser
	tempDir string
	log     *log.Helper
}

func NewChromePool(logger log.Logger) (*ChromePool, func()) {
	dataStore, err := os.MkdirTemp("", "btl-*")
	if err != nil {
		panic(err)
	}
	chromeLauncher := launcher.New().
		NoSandbox(true).
		Headless(true).
		Leakless(false).
		Set("disable-gpu", "true").
		Set("ignore-certificate-errors", "true").
		Set("ignore-certificate-errors", "1").
		Set("disable-crash-reporter", "true").
		Set("disable-notifications", "true").
		Set("hide-scrollbars", "true").
		Set("window-size", fmt.Sprintf("%d,%d", 1080, 1920)).
		Set("mute-audio", "true").
		Delete("use-mock-keychain").
		Env(append(os.Environ(), "TZ=Asia/Shanghai")...).
		UserDataDir(dataStore)

	if runtime.GOOS == "darwin" {
		chromeLauncher.Headless(false)
	}

	p, installed := launcher.LookPath()
	if !installed {
		panic("chrome not installed")
	}
	lUrl, err := chromeLauncher.Bin(p).Launch()
	if err != nil {
		panic(err)
	}
	browser := rod.New().ControlURL(lUrl)
	if browserErr := browser.Connect(); browserErr != nil {
		panic(browserErr)
	}

	cancel := func() {
		_ = browser.Close()
		_ = os.RemoveAll(dataStore)
	}

	return &ChromePool{
		log:     log.NewHelper(logger),
		browser: browser,
	}, cancel
}

func (c ChromePool) NavigateRequest(url string) (body string, err error) {
	page, err := c.browser.Page(proto.TargetCreateTarget{
		URL: url,
	})
	if err != nil {
		return
	}
	defer page.Close()

	//https://github.com/go-rod/rod/issues/230
	err = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent:      UserAgent,
		AcceptLanguage: Language,
	})
	if err != nil {
		return
	}

	_ = page.WaitLoad()

	ctxTimeout, cancel := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancel()

	time.Sleep(5 * time.Second)
	body, err = page.Context(ctxTimeout).HTML()
	return
}

func (c ChromePool) ScreenShot(url string, imgPath string) error {
	c.log.Infof("start ScreenShot for url <%s>", url)
	page, err := c.browser.Page(proto.TargetCreateTarget{
		URL: url,
	})
	if err != nil {
		return err
	}
	defer page.Close()

	//https://github.com/go-rod/rod/issues/230
	_ = page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent:      UserAgent,
		AcceptLanguage: Language,
	})
	_ = page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:  1440,
		Height: 900,
	})

	_ = page.WaitLoad()

	time.Sleep(10 * time.Second)

	ctxTimeout, cancel := context.WithTimeout(context.TODO(), time.Minute*2)
	defer cancel()

	quality := 100
	imgData, err := page.Context(ctxTimeout).Screenshot(false, &proto.PageCaptureScreenshot{
		Quality: &quality,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(imgPath, imgData, os.ModePerm)
}
