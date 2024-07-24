package strava

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
)

type UploadFit struct {
	fit       string
	ctx       context.Context
	log       *log.Helper
	uploadUrl string
	username  string
	password  string
}

func NewUploadFit(fit string, logger log.Logger) *UploadFit {
	return &UploadFit{
		fit:       fit,
		ctx:       context.Background(),
		log:       log.NewHelper(logger),
		uploadUrl: "https://www.strava.com/upload/select",
		username:  "beer5214@126.com",
		password:  "Hm%L-sjRV7RaGG+",
	}
}

func (s *UploadFit) Run() error {
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(s.ctx, allocOpts...)
	defer cancel()

	ctx, _ := chromedp.NewContext(allocCtx)

	err := s.login(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.uploadFile(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *UploadFit) uploadFile(ctx context.Context) error {
	uploadXpath := `//*[@id="uploadFile"]/form/input[3]`
	actions := chromedp.Tasks{
		chromedp.WaitVisible(uploadXpath, chromedp.BySearch),
		chromedp.SendKeys(uploadXpath, s.fit, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(10 * time.Second)
			return nil
		}),

		//todo
	}
	return chromedp.Run(ctx, actions)
}

func (s *UploadFit) login(ctx context.Context) error {
	err := chromedp.Run(ctx,
		chromedp.Navigate(s.uploadUrl),
		chromedp.WaitVisible(`//*[@id="login_form"]`, chromedp.BySearch),
		chromedp.SendKeys(`//*[@id="email"]`, s.username, chromedp.BySearch),
		chromedp.SendKeys(`//*[@id="password"]`, s.password, chromedp.BySearch),
		chromedp.Click(`//*[@id="login-button"]`, chromedp.BySearch),

		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(5 * time.Second)

			var (
				currentURL string
				err        error
			)
			err = chromedp.Location(&currentURL).Do(ctx)
			if err != nil {
				return err
			}
			if currentURL == s.uploadUrl {
				s.log.Info("login success")
			} else {
				return errors.New("login failed")
			}
			return nil
		}),
	)
	if err != nil {
		s.log.Errorf("chrome run err %+v", err)
		return err
	}
	return nil
}
