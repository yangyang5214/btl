package pkg

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type Json2Csv struct {
	jsonFile string
}

func NewJson2Csv(jsonFile string) *Json2Csv {
	return &Json2Csv{jsonFile: jsonFile}
}

func (s *Json2Csv) Run() error {
	bytes, err := os.ReadFile(s.jsonFile)
	if err != nil {
		log.Errorf("read file error %v", err)
		return errors.WithStack(err)
	}
	content := string(bytes)
	var datas []map[string]interface{}
	if strings.HasPrefix(content, "[") {
		if err = json.NewDecoder(strings.NewReader(content)).Decode(&datas); err != nil {
			return errors.WithStack(err)
		}
	} else {
		lines := strings.Split(content, "\n")
		if strings.HasPrefix(lines[0], "{") && strings.HasSuffix(lines[0], "}") {
			for _, line := range lines {
				if line == "" {
					break
				}
				var data map[string]interface{}
				err = json.Unmarshal([]byte(line), &data)
				if err != nil {
					log.Errorf("unmarshal json %s error %v", line, err)
					return errors.WithStack(err)
				}
				datas = append(datas, data)
			}
		} else {
			//todo
			return errors.New("file format not supported")
		}
	}
	return saveCsv(getFilePrefix(s.jsonFile)+".csv", datas)
}
