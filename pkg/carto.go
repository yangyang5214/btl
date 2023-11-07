package pkg

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg/utils"
	sm "github.com/yangyang5214/go-staticmaps"
)

var (
	mapName = map[string][]string{
		"carto-light": {"https://cartodb-basemaps-b.global.ssl.fastly.net/light_all/%d/%d/%d.png"},
		"carto-dark":  {"https://cartodb-basemaps-b.global.ssl.fastly.net/dark_all/%d/%d/%d.png"},
		"osm": {
			"https://a.tile.openstreetmap.org/%d/%d/%d.png",
			"https://b.tile.openstreetmap.org/%d/%d/%d.png",
			"https://c.tile.openstreetmap.org/%d/%d/%d.png",
		},
	}
)

type Carto struct {
	name      string
	smCtx     *sm.Context
	zoom      int
	client    *http.Client
	cacheDir  string
	userAgent string

	center   s2.LatLng
	radiusKm float64
}

func NewCarto(name string, zoom int, center s2.LatLng, radiusKm float64) *Carto {
	return &Carto{
		name:      name,
		smCtx:     sm.NewContext(),
		zoom:      zoom,
		client:    http.DefaultClient,
		center:    center,
		cacheDir:  os.Getenv("CARTO_CACHE"),
		userAgent: "Mozilla/5.0+(compatible; go-staticmaps/0.1; https://github.com/flopp/go-staticmaps)",
		radiusKm:  radiusKm,
	}
}

func (c *Carto) Run() (err error) {
	if c.cacheDir == "" {
		return errors.New("cache dir is required")
	}
	urls, ok := mapName[c.name]
	if !ok {
		return nil
	}

	deltaLat := c.radiusKm / 111.0
	deltaLng := c.radiusKm / 111.0 * math.Cos(float64(c.center.Lat*math.Pi/180.0))

	start := s2.LatLngFromDegrees(c.center.Lat.Degrees()+deltaLat, c.center.Lng.Degrees()-deltaLng)
	end := s2.LatLngFromDegrees(c.center.Lat.Degrees()-deltaLat, c.center.Lng.Degrees()+deltaLng)

	log.Infof("zoom %d, start: %v,%v, end: %v,%v", c.zoom, start.Lng.Degrees(), start.Lat.Degrees(), end.Lng.Degrees(), end.Lat.Degrees())
	bounds, err := utils.GenBounds(start, end, c.zoom)
	if err != nil {
		return err
	}

	log.Infof("bounds is %v", bounds)

	var wg sync.WaitGroup
	var finalUrl string
	workerCh := make(chan struct{}, 100)
	for i := bounds.X[0]; i < bounds.X[1]; i++ {
		for j := bounds.Y[0]; j < bounds.Y[1]; j++ {
			finalUrl = urls[0]
			if len(urls) != 1 {
				s := rand.NewSource(time.Now().Unix())
				r := rand.New(s)
				finalUrl = urls[r.Intn(len(urls))]
			}

			wg.Add(1)
			workerCh <- struct{}{}
			go func(x, y int, urlStr string) {
				defer func() {
					wg.Done()
					<-workerCh
				}()

				for k := 0; k < 3; k++ { //重试三次
					err = c.download(x, y, urlStr)
					if err != nil {
						log.Errorf("error download %v", err)
					} else {
						break
					}
				}
			}(i, j, finalUrl)
		}
	}
	wg.Wait()

	return nil
}

func (c *Carto) download(x, y int, urlStr string) error {
	fileName := path.Join(
		c.cacheDir,
		c.name,
		strconv.Itoa(c.zoom),
		strconv.Itoa(x),
		strconv.Itoa(y),
	)

	//exist
	_, err := os.Stat(fileName)
	if err == nil {
		return nil
	}

	finalUrl := fmt.Sprintf(urlStr, c.zoom, x, y)
	log.Infof("start download %s", finalUrl)
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", c.userAgent)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("response code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = c.storeCache(fileName, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Carto) createCacheDir(path string) error {
	src, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, 0755)
		}
		return err
	}
	if src.IsDir() {
		return nil
	}

	return fmt.Errorf("file exists but is not a directory: %s", path)
}

func (c *Carto) storeCache(fileName string, data []byte) error {
	dir, _ := filepath.Split(fileName)

	if err := c.createCacheDir(dir); err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
