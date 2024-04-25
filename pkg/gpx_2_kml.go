package pkg

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/twpayne/go-kml"
	"golang.org/x/image/colornames"
	"os"
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
	op := new(options)
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
	points, err := s.getAllPoints()

	highlightStyle := kml.SharedStyle(
		"highlight",
		kml.IconStyle(
			kml.Icon(
				kml.Href("http://maps.google.com/mapfiles/kml/paddle/red-stars.png"),
			),
		),
		kml.LineStyle(
			kml.Color(colornames.Red),
			kml.Width(8),
		),
	)
	normalStyle := kml.SharedStyle(
		"normal",
		kml.IconStyle(
			kml.Icon(
				kml.Href("http://maps.google.com/mapfiles/kml/paddle/wht-blank.png"),
			),
		),
		kml.LineStyle(
			kml.Color(colornames.Blue),
			kml.Width(6),
		),
	)

	styleMap := kml.SharedStyleMap(
		"styleMap",
		kml.Pair(
			kml.Key(StyleStateNormal),
			kml.StyleURL(normalStyle.URL()),
		),
		kml.Pair(
			kml.Key(StyleStateHighlight),
			kml.StyleURL(highlightStyle.URL()),
		),
	)

	k := kml.KML(
		kml.Document(
			kml.StyleMap(styleMap),
			kml.Style(normalStyle),
			kml.Style(highlightStyle),
			kml.Placemark(
				kml.Name(s.opts.name),
				kml.GxTrack(points...),
				kml.StyleURL(styleMap.URL()),
			),
		),
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
