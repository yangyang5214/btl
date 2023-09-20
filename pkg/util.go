package pkg

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

func downloadSave(client *http.Client, urlStr string, filename string) error {
	resp, err := client.Get(urlStr)
	if err != nil {
		log.Errorf("http request failed %+v. %s", err, urlStr)
		return err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	f, err := os.Create(filename)
	if err != nil {
		log.Errorf("create file %v error", filename)
		return err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
