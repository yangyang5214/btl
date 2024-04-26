package pkg

import (
	"encoding/xml"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/twpayne/go-kml"
	"golang.org/x/image/colornames"
	"os"
	"path/filepath"
)

type options struct {
	name       string
	resultFile string
}

type Option func(*options)

func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

func WithResultFile(resultFile string) Option {
	return func(o *options) { o.resultFile = resultFile }
}

type Gpx2Kml struct {
	gpxFile string
	log     *log.Helper
	opts    options
}

// NewGpx2Kml
// https://developers.google.com/kml/documentation/topicsinkml?hl=zh-cn
// https://github.com/twpayne/go-kml
func NewGpx2Kml(gpxFile string, logger log.Logger, opts ...Option) *Gpx2Kml {
	op := &options{
		name:       "轨迹",
		resultFile: "result.kml",
	}
	for _, opt := range opts {
		opt(op)
	}
	return &Gpx2Kml{
		gpxFile: gpxFile,
		log:     log.NewHelper(logger),
		opts:    *op,
	}
}

func (s *Gpx2Kml) getAllPoints() (points []kml.Element, err error) {
	gpxData, err := gpx.ParseFile(s.gpxFile)
	if err != nil {
		return
	}
	for _, track := range gpxData.Tracks {
		for _, segment := range track.Segments {
			for _, p := range segment.Points {
				points = append(points, kml.When(p.Timestamp))
			}
		}
	}

	for _, track := range gpxData.Tracks {
		for _, segment := range track.Segments {
			for _, p := range segment.Points {
				points = append(points, kml.GxCoord(kml.Coordinate{
					Lon: p.Longitude,
					Lat: p.Latitude,
					Alt: p.Elevation.Value(),
				}))
			}
		}
	}
	return
}

const (
	StyleStateNormal    kml.StyleStateEnum = "normal"
	StyleStateHighlight kml.StyleStateEnum = "highlight"
)

func (s *Gpx2Kml) Run() error {
	s.log.Infof("import gpx file %s", s.gpxFile)
	p, err := filepath.Abs(s.gpxFile)
	if err != nil {
		return err
	}
	s.gpxFile = p

	points, err := s.getAllPoints()

	normalStyle := kml.SharedStyle(
		"multiTrack_n",
		kml.IconStyle(
			kml.Icon(
				kml.Href("https://earth.google.com/images/kml-icons/track-directional/track-0.png"),
			),
		),
		kml.LineStyle(
			kml.Color(colornames.Red),
			kml.Width(3),
		),
	)

	multiTrack := kml.SharedStyleMap(
		string(StyleStateNormal),
		kml.Pair(
			kml.Key(StyleStateNormal),
			kml.StyleURL(normalStyle.URL()),
		),
	)

	k := kml.KML(
		kml.Document(
			kml.Name(s.opts.name+".kml"),
			multiTrack,
			normalStyle,
			kml.Placemark(
				kml.Name(s.opts.name),
				kml.StyleURL(multiTrack.URL()),
				kml.GxTrack(points...),
			),
		),
	)

	k.Attr = append(k.Attr,
		xml.Attr{Name: xml.Name{Local: "xmlns:gx"}, Value: "http://www.google.com/kml/ext/2.2"},
		xml.Attr{Name: xml.Name{Local: "xmlns:kml"}, Value: "http://www.opengis.net/kml/2.2"},
		xml.Attr{Name: xml.Name{Local: "xmlns:atom"}, Value: "http://www.w3.org/2005/Atom"},
	)

	f, err := os.Create(s.opts.resultFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	if err := k.WriteIndent(f, "", "  "); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
