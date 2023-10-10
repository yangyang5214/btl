package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Markdown struct {
	filepath   string
	httpClient *http.Client
}

func NewMarkdown(filepath string) *Markdown {
	return &Markdown{
		filepath: filepath,
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (m *Markdown) downloadImag(line string) (string, error) {
	var err error

	var imgUrl string
	if strings.HasPrefix(line, "![") {
		imgUrl = m.parserMarkdownImgUrl(line)
	} else if strings.HasPrefix(line, "<img") {
		imgUrl = parserSrcImgUrl(line)
	} else {
		return line, nil
	}

	pwd, err := GetPwd()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(imgUrl, "http") {
		imgType := m.parserImageType(imgUrl)
		if imgType == "" {
			return "", errors.New(fmt.Sprintf("imgUrl %s not expected", imgUrl))
		}
		filename := md5Sum(line) + "." + imgType

		err = downloadSave(m.httpClient, imgUrl, path.Join(pwd, "images", filename))
		if err != nil {
			return "", err
		}
		line = fmt.Sprintf("![](./images/%s)", filename)
	}
	return line, err
}

// ParseImages is 解析图片链接，转换为本地文件依赖
func (m *Markdown) ParseImages() error {
	bytes, err := os.ReadFile(m.filepath)
	if err != nil {
		log.Errorf("open file %s error: %s", m.filepath, err)
		return err
	}
	lines := strings.Split(string(bytes), "\n")

	var result []string
	for _, line := range lines {
		line, err = m.downloadImag(line)
		if err != nil {
			return err
		}
		result = append(result, line)
	}

	err = os.Remove(m.filepath)
	if err != nil {
		return err
	}

	err = os.WriteFile(m.filepath, []byte(strings.Join(result, "\n")), 0755)
	if err != nil {
		return err
	}

	return nil
}

func (m *Markdown) parserImageType(urlStr string) string {
	length := len(urlStr)
	for i := length - 1; i > 0; i-- {
		if string(urlStr[i]) == "." {
			return urlStr[i+1 : length]
		}
	}
	return ""
}

func parserSrcImgUrl(content string) string {
	//skip ">
	end := len(content) - 3
	start := 0
	for i := end; i > 0; i-- {
		if string(content[i]) == `"` {
			start = i
			break
		}
	}
	return strings.Trim(content[start:end+1], "\"")
}
func (m *Markdown) parserMarkdownImgUrl(content string) string {
	var start, end int
	for index, item := range content {
		if string(item) == "(" {
			start = index + 1
		} else if string(item) == ")" {
			end = index
		}
	}
	return content[start:end]
}
