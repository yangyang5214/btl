package pkg

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
)

func downloadSave(client *http.Client, urlStr string, filepath string) error {
	if _, err := os.Stat(filepath); err == nil {
		log.Infof("file %s exists skip", filepath)
		return nil
	}

	resp, err := client.Get(urlStr)
	if err != nil {
		log.Errorf("http request failed %+v. %s", err, urlStr)
		return err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)

	err = os.MkdirAll(path.Dir(filepath), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath)
	if err != nil {
		log.Errorf("create file %v error", filepath)
		return err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func md5Sum(text string) string {
	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func readJsonLineFile(p string) ([]map[string]interface{}, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var result []map[string]interface{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		var data map[string]interface{}
		err = json.Unmarshal([]byte(line), &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	return result, nil
}
