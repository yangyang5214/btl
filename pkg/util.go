package pkg

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

func downloadSave(client *http.Client, urlStr string, filepath string) error {
	if _, err := os.Stat(filepath); err == nil {
		log.Infof("file %s exists skip", filepath)
		return nil
	}

	resp, err := client.Get(urlStr)
	if err != nil {
		log.Errorf("http request failed %+v. %s", err, urlStr)
		return errors.WithStack(err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)

	err = os.MkdirAll(path.Dir(filepath), 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.Create(filepath)
	if err != nil {
		log.Errorf("create file %v error", filepath)
		return errors.WithStack(err)
	}
	_, err = f.Write(bytes)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func md5Sum(text string) string {
	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func readCsv(p string) ([]map[string]interface{}, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()
	var result []map[string]interface{}

	reader := csv.NewReader(f)
	rawData, err := reader.ReadAll()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var headers []string
	for lineNum, record := range rawData {
		if lineNum == 0 {
			for i := 0; i < len(record); i++ {
				headers = append(headers, strings.TrimSpace(record[i]))
			}
		} else {
			line := map[string]interface{}{}
			for i := 0; i < len(record); i++ {
				line[headers[i]] = record[i]
			}
			result = append(result, line)
		}
	}
	return result, nil
}

func readJsonLineFile(p string) ([]map[string]interface{}, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()
	var result []map[string]interface{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		var data map[string]interface{}
		err = json.Unmarshal([]byte(line), &data)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result = append(result, data)
	}
	return result, nil
}

func getFilePrefix(p string) string {
	paths := strings.Split(p, "/")
	last := paths[len(paths)-1]
	a := strings.Split(last, ".")
	return a[0]
}

func mapKeys(m map[string]interface{}) []string {
	var r []string
	for k := range m {
		r = append(r, k)
	}
	return r
}

func saveCsv(filename string, datas []map[string]interface{}) error {
	outputFile, err := os.Create(filename)
	if err != nil {
		log.Errorf("create file error: %+v", err)
		return errors.WithStack(err)
	}
	defer outputFile.Close()
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()
	headers := mapKeys(datas[0])
	if err = writer.Write(headers); err != nil {
		log.Errorf("write headers error: %+v", err)
		return errors.WithStack(err)
	}
	for _, r := range datas {
		var csvRow []string

		for _, header := range headers {
			v, _ := r[header]
			csvRow = append(csvRow, fmt.Sprintf("%v", v))
		}

		if err = writer.Write(csvRow); err != nil {
			log.Errorf("write csv row data error: %+v", err)
			return errors.WithStack(err)
		}
	}
	return nil
}

func ColorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02X%02X%02X", r>>8, g>>8, b>>8)
}

func FileExists(fp string) bool {
	if fp == "" {
		return false
	}
	_, err := os.Stat(fp)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
