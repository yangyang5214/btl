package pkg

type Gpx2Kml struct {
	gpxFile string
}

// NewGpx2Kml
// https://developers.google.com/kml/documentation/topicsinkml?hl=zh-cn
// https://github.com/twpayne/go-kml
func NewGpx2Kml(gpxFile string) *Gpx2Kml {
	return &Gpx2Kml{
		gpxFile: gpxFile,
	}
}
