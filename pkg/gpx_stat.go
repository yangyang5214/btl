package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
	"math"
	"time"
)

type GpxStat struct {
	file string
}

type GpxStatModel struct {
	TimeDuration string
	Speed        float64 // 速度 km/h
	Elevation    int     // 海拔
}

func (gsm *GpxStatModel) String() string {
	return fmt.Sprintf("%v -> %v -> %v", gsm.TimeDuration, gsm.Speed, gsm.Elevation)
}

func NewGpxStat(f string) *GpxStat {
	return &GpxStat{file: f}
}

func (g *GpxStat) Run() ([]*GpxStatModel, error) {
	gpxData, err := gpx.ParseFile(g.file)
	if err != nil {
		return nil, err
	}
	points := gpxData.Tracks[0].Segments[0].Points
	startTime := points[0].Timestamp
	log.Infof("point length is: %d", len(points))

	var r []*GpxStatModel
	for i := 1; i < len(points); i += 50 {
		speed := calculateSpeed(points[i-1], points[i])
		subTime := points[i].Timestamp.Sub(startTime)
		elv := points[i].Elevation.Value()

		r = append(r, &GpxStatModel{
			TimeDuration: formatSecond(subTime),
			Speed:        speed * 3.6,
			Elevation:    int(elv),
		})
	}
	return r, nil
}

func formatSecond(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func calculateSpeed(point1, point2 gpx.GPXPoint) float64 {
	distance := haversine(point1.Latitude, point1.Longitude, point2.Latitude, point2.Longitude)
	timeDiff := point2.Timestamp.Sub(point1.Timestamp).Seconds()
	return distance / timeDiff
}

// Haversine formula to calculate distance between two points on the Earth
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // in meters

	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)

	lat1 = lat1 * (math.Pi / 180)
	lat2 = lat2 * (math.Pi / 180)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return distance
}
