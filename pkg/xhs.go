package pkg

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/image/webp"
)

type Xhs struct {
	urlStr     string
	httpClient *http.Client
	log        *log.Helper
}

func NewXhs(u string) *Xhs {
	return &Xhs{
		urlStr:     u,
		httpClient: http.DefaultClient,
		log:        log.NewHelper(log.DefaultLogger),
	}
}

func (s *Xhs) Run() error {
	if s.urlStr == "" {
		return nil
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout to prevent hanging
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	headers := map[string]interface{}{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
	}
	if err := chromedp.Run(ctx, navigate(s.urlStr, headers)); err != nil {
		return err
	}

	var htmlContent string
	if err := chromedp.Run(ctx, extractHTML(&htmlContent)); err != nil {
		return err
	}

	//save index.html
	err := s.writeFile(htmlContent, "index.html")
	if err != nil {
		return err
	}
	r, err := xpathUrl(htmlContent)
	if err != nil {
		return err
	}
	return s.saveDir(r)
}

func (s *Xhs) saveDir(r *Result) error {
	arrs := strings.Split(s.urlStr, "/")
	key := arrs[len(arrs)-1]
	_ = os.MkdirAll(key, 0755)
	readme := path.Join(key, "README.md")
	err := s.writeFile(r.Content, readme)
	if err != nil {
		return err
	}

	if len(r.ImageUrls) > 5 {
		r.ImageUrls = r.ImageUrls[:5]
	}
	for index, urlStr := range r.ImageUrls {
		p := path.Join(key, fmt.Sprintf("%d.jpg", index))
		s.log.Infof("save img %s", p)
		imgData, err := s.download(urlStr)
		if err != nil {
			return err
		}
		err = s.saveWebpImg(imgData, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Xhs) writeFile(content string, p string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, _ = f.WriteString(content)
	return nil
}

func (s *Xhs) saveWebpImg(data []byte, p string) error {
	webpImage, err := webp.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}

	jpegFile, err := os.Create(p)
	if err != nil {
		return err
	}
	defer jpegFile.Close()

	// Encode the WebP image to JPEG and save it to the file
	err = jpeg.Encode(jpegFile, webpImage, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}
	return nil
}

func (s *Xhs) download(u string) ([]byte, error) {
	resp, err := s.httpClient.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func navigate(url string, headers map[string]interface{}) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(headers),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
	}
}

func extractHTML(target *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			root, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			node, err := dom.GetOuterHTML().WithNodeID(root.NodeID).Do(ctx)
			if err != nil {
				return err
			}
			*target = node
			return nil
		}),
	}
}

type Result struct {
	ImageUrls []string
	Content   string
}

func xpathUrl(html string) (*Result, error) {
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	list, err := htmlquery.QueryAll(doc, "//meta[@name='og:image']//@content")
	if err != nil {
		return nil, err
	}

	m := make(map[string]bool)
	var result []string
	for _, node := range list {
		v := node.FirstChild.Data
		_, ok := m[v]
		if !ok {
			m[v] = true
			result = append(result, v)
		}
	}
	content, err := htmlquery.Query(doc, "//meta[@name='description']/@content")
	if err != nil {
		return nil, err
	}
	return &Result{
		ImageUrls: result,
		Content:   content.FirstChild.Data,
	}, nil
}
