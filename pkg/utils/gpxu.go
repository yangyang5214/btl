package utils

import (
	"fmt"
	. "image/color"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/qichengzx/coordtransform"

	"github.com/yangyang5214/btl/pkg/model"

	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
)

func ParseGpxData(files []string) ([]*gpx.GPX, error) {
	var results []*gpx.GPX
	for _, f := range files {
		p, err := filepath.Abs(f)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		gpxData, err := gpx.ParseFile(p)
		if err != nil {
			log.Errorf("gpx parse file <%s> error: %v", p, err)
			return nil, errors.WithStack(err)
		}
		if len(gpxData.Tracks) == 0 {
			continue
		}
		if gpxData.GetTrackPointsNo() == 0 {
			continue
		}
		results = append(results, gpxData)
	}
	return results, nil
}

func ParsePositions(datas []*gpx.GPX) ([][]s2.LatLng, error) {
	var positions [][]s2.LatLng
	for _, gpxData := range datas {
		var local []s2.LatLng
		for _, trk := range gpxData.Tracks {
			for _, seg := range trk.Segments {
				for _, pt := range seg.Points {
					local = append(local, s2.LatLngFromDegrees(pt.GetLatitude(), pt.GetLongitude()))
				}
			}
		}
		positions = append(positions, local)
	}
	return positions, nil
}

func CountPoints(positions [][]s2.LatLng) int {
	var r int
	for _, sub := range positions {
		r = r + len(sub)
	}
	return r
}

func FindGpxFiles(dirPath string) []string {
	if strings.HasSuffix(dirPath, ".gpx") {
		return []string{dirPath}
	}
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil
	}
	d, err := os.ReadDir(absPath)
	if err != nil {
		return nil
	}
	var r []string
	for _, item := range d {
		name := item.Name()
		p := path.Join(absPath, name)
		if item.IsDir() {
			r = append(r, FindGpxFiles(p)...)
		}
		if strings.HasSuffix(name, ".gpx") {
			r = append(r, p)
		}
	}
	return r
}

func FindFitFiles(dirPath string) []string {
	if strings.HasSuffix(dirPath, ".fit") {
		return []string{dirPath}
	}
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil
	}
	d, err := os.ReadDir(absPath)
	if err != nil {
		return nil
	}
	var r []string
	for _, item := range d {
		name := item.Name()
		p := path.Join(absPath, name)
		if item.IsDir() {
			r = append(r, FindGpxFiles(p)...)
		}
		if strings.HasSuffix(name, ".fit") {
			r = append(r, p)
		}
	}
	return r
}

func GetColor(index int, colors []Color) Color {
	return colors[index%len(colors)]
}

func LonLatToTileCoordinates(lng, lat float64, zoom int) (x, y int) {
	x = int(math.Floor((lng + 180.0) / 360.0 * math.Pow(2.0, float64(zoom))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * math.Pow(2.0, float64(zoom))))
	return x, y
}

func ParserLatLng(location string) (s2.LatLng, error) {
	point := strings.Split(location, ",")
	if len(point) != 2 {
		return s2.LatLng{}, errors.New("Invalid carto location")
	}
	lat, err := strconv.ParseFloat(point[1], 64)
	if err != nil {
		return s2.LatLng{}, err
	}
	lng, err := strconv.ParseFloat(point[0], 64)
	if err != nil {
		return s2.LatLng{}, err
	}
	return s2.LatLngFromDegrees(lat, lng), nil
}

func GenBounds(start s2.LatLng, end s2.LatLng, zoom int) (*model.Bounds, error) {
	startX, startY := LonLatToTileCoordinates(start.Lng.Degrees(), start.Lat.Degrees(), zoom)
	endX, endY := LonLatToTileCoordinates(end.Lng.Degrees(), end.Lat.Degrees(), zoom)
	return &model.Bounds{
		X: []int{startX, endX},
		Y: []int{startY, endY},
	}, nil
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (l LatLng) String() string {
	return fmt.Sprintf("%f,%f", l.Lng, l.Lat)
}

// GCJ02String https://lbs.amap.com/api/javascript-api-v2/guide/transform/convertfrom
func (l LatLng) GCJ02String() string {
	lng, lat := coordtransform.WGS84toGCJ02(l.Lng, l.Lat)
	return fmt.Sprintf("%f,%f", lng, lat)
}

func GetPoints(gpxFiles []string) ([][]*LatLng, error) {
	var gpxs []*gpx.GPX
	for _, gf := range gpxFiles {
		gpxData, err := gpx.ParseFile(gf)
		if err != nil {
			return nil, err
		}
		gpxs = append(gpxs, gpxData)
	}
	return GetPointsFromGpx(gpxs)
}

func SortGpx(gpxs []*gpx.GPX) []*gpx.GPX {
	dateMap := make(map[int64]*gpx.GPX)
	keys := make([]int64, 0, len(dateMap))
	for _, f := range gpxs {
		startTime := f.TimeBounds().StartTime.UnixMilli()
		keys = append(keys, startTime)
		dateMap[startTime] = f
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	var result []*gpx.GPX
	for _, k := range keys {
		result = append(result, dateMap[k])
	}
	return result
}

// GetPointsFromGpx 外部调用
func GetPointsFromGpx(gpxs []*gpx.GPX) ([][]*LatLng, error) {
	gpxs = SortGpx(gpxs)
	var result [][]*LatLng
	for _, gpxData := range gpxs {
		var points []*LatLng
		for _, track := range gpxData.Tracks {
			for _, seg := range track.Segments {
				for _, pt := range seg.Points {
					points = append(points, &LatLng{
						Lat: pt.GetLatitude(),
						Lng: pt.GetLongitude(),
					})
				}
			}
		}
		result = append(result, points)
	}
	return result, nil
}
