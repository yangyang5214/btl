package pkg

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"strings"

	"net/http"
	"os"
	"path"
	"time"
)

type Github struct {
	token string

	httpClient *http.Client
}

func readToken() (string, error) {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	f, err := os.ReadFile(path.Join(dirName, ".github_token"))
	if err != nil {
		return "", err
	}
	return string(f), nil
}

func NewGithub() *Github {
	token, err := readToken()
	if err != nil {
		log.Fatal("read github token failed. from $HOME/.github_token")
	}
	return &Github{
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
		token: strings.Trim(token, "\n"),
	}
}

func (g *Github) DownloadGist(gistId string) (string, error) {
	apiUrl := "https://api.github.com/gists/" + gistId
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Errorf("do http error: %v", err)
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("X-Github-Api-Version", "2022-11-28")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	data := gjson.ParseBytes(bytes)
	files := data.Get("files").Map()
	for _, result := range files {
		filename := result.Get("filename").String()
		log.Infof("start downloading file %s", filename)
		rawUrl := result.Get("raw_url").String()
		err = downloadSave(g.httpClient, rawUrl, filename)
		if err != nil {
			return "", err
		}
		return filename, nil
	}
	return "", nil
}
