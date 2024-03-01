package heatmap

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
)

type MergeGpx struct {
	httpClient *http.Client
	headers    *http.Header
}

func NewMergeGpx() *MergeGpx {
	return &MergeGpx{
		httpClient: http.DefaultClient,
	}
}

func (m *MergeGpx) initReq(urlStr string, body string) (*http.Request, error) {
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic bG9uZzp6aGFuZzUyMTQ/")
	req.Header.Set("token", "621oZth46_KPrWRQMg_giK2pecHL41s7be")
	return req, nil
}

func (m *MergeGpx) GetActivityIds() ([]string, error) {
	url := "https://api.beer5214.com/heatmap/ids"
	body := make(map[string]interface{})
	body["merged"] = true
	body["deleted"] = false
	body["activity_type"] = 2
	bodyStr, _ := json.Marshal(body)

	req, err := m.initReq(url, string(bodyStr))
	resp, err := m.httpClient.Do(req)
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := gjson.ParseBytes(bytes)
	var ids []string
	for _, result := range data.Get("ids").Array() {
		ids = append(ids, result.String())
	}
	return ids, nil
}

func (m *MergeGpx) DownloadActivity(activityId string) ([]byte, error) {
	url := "https://api.beer5214.com/heatmap/line"
	body := make(map[string]interface{})
	body["activity_id"] = activityId
	bodyStr, _ := json.Marshal(body)

	req, err := m.initReq(url, string(bodyStr))
	resp, err := m.httpClient.Do(req)
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
