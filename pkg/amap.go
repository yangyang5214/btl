package pkg

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type Amap struct {
	APIKey string
	log    *log.Helper
}

func NewAmap(log *log.Helper) *Amap {
	homeDir, _ := os.UserHomeDir()
	bytes, err := os.ReadFile(path.Join(homeDir, ".amap_key"))
	if err != nil {
		panic(err)
	}

	return &Amap{
		APIKey: strings.Trim(string(bytes), "\n"),
		log:    log,
	}
}

func (a *Amap) GetLocationByAddress(address string) (string, error) {
	urlStr := "https://restapi.amap.com/v3/geocode/geo"
	params := url.Values{}
	params.Set("address", address)
	params.Set("key", a.APIKey)
	params.Set("output", "json")

	urlStr = urlStr + "?" + params.Encode()
	a.log.Infof("amap url: %s", urlStr)

	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	gjsonResult := gjson.ParseBytes(body)
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("amap resp status is: " + resp.Status)
	}

	status := gjsonResult.Get("status").Int()
	if status != 1 {
		errInfo := gjsonResult.Get("info").String()
		a.log.Errorf(errInfo)
		return "", errors.New("amap resp status is not 1")
	}

	locations := gjsonResult.Get("geocodes.0.location").String()
	if locations == "" {
		return "", errors.New("amap resp location is empty")
	}
	//Longitude,Latitude 经度/纬度
	return locations, nil
}
