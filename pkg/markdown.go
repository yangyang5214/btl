package pkg

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Markdown struct {
	filename string
}

func NewMarkdown(filename string) *Markdown {
	return &Markdown{
		filename: filename,
	}
}

// parseImages is 解析图片链接，转换为本地文件依赖
func (m *Markdown) parseImages() error {
	bytes, err := os.ReadFile(m.filename)
	if err != nil {
		log.Errorf("open file %s error: %s", m.filename, err)
		return err
	}
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "![") {
			line = "111"
		}
	}
	return nil
}
