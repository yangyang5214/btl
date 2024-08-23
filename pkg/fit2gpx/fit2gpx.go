package fit2gpx

import (
	_ "embed"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
)

type Fit2Gpx struct {
	fitFile    string
	log        *log.Helper
	workDir    string
	fitCSVTool string
	resultGpx  string
}

func NewFit2Gpx(fitFile string, logger log.Logger) *Fit2Gpx {
	homeDir, _ := os.UserHomeDir()
	workDir := path.Join(homeDir, ".fit2gpx")
	_ = os.MkdirAll(workDir, 0755)
	return &Fit2Gpx{
		fitFile:    fitFile,
		log:        log.NewHelper(logger),
		workDir:    workDir,
		fitCSVTool: path.Join(workDir, "FitCSVTool.jar"),
		resultGpx:  "result.gpx",
	}
}

func (s *Fit2Gpx) SetResultPath(resultPath string) {
	s.resultGpx = resultPath
}

type Session struct {
	points []*Record

	sportType string
}

type Record struct {
	Ts        int64
	Lat       float64
	Lng       float64
	Altitude  float64
	HeartRate int64
	Distance  float64
	Speed     float64
	Cadence   int64
}

func (s *Fit2Gpx) fit2Csv(csvPath string) error {
	cmd := fmt.Sprintf("java -jar -Xmx2g %s -b %s %s", s.fitCSVTool, s.fitFile, csvPath)
	s.log.Infof("run cmd %s", cmd)
	return exec.Command("/bin/bash", "-c", cmd).Run()
}

func (s *Fit2Gpx) parserCsv(csvPath string) (*Session, error) {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	session := &Session{}
	r := csv.NewReader(csvFile)
	r.FieldsPerRecord = -1

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if record[2] == "sport" && record[0] == "Data" {
			for i := 3; i < len(record); i += 3 {
				if record[i] == "sport" {
					session.sportType = ParseSport(record[i+1])
					break
				}
			}
		}

		if record[2] == "record" {
			robj, err := s.parserRecord(record)
			if err != nil {
				return nil, err
			}
			if robj != nil {
				session.points = append(session.points, robj)
			}
		}
	}
	return session, err
}

func (s *Fit2Gpx) parserRecord(msgs []string) (*Record, error) {
	m := make(map[string]string)
	for i := 3; i < len(msgs)-3; i += 3 {
		k := msgs[i]
		v := msgs[i+1]
		if k == "" && v == "" {
			break
		}
		m[msgs[i]] = msgs[i+1]
	}

	ts, ok := m["timestamp"]
	if !ok {
		return nil, nil
	}
	if ts == "1" {
		return nil, nil //ignore
	}

	tsInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil, nil
	}

	lat := parseFloat(m["position_lat"])
	if lat == 0 {
		return nil, nil //ignore
	}

	return &Record{
		Ts:        tsInt + 631036800,
		Lat:       lat,
		Lng:       parseFloat(m["position_long"]),
		Altitude:  parserAttr(m, []string{"altitude", "enhanced_altitude"}),
		HeartRate: parseInt(m["heart_rate"]),
		Distance:  parseFloat(m["distance"]),
		Speed:     parserAttr(m, []string{"speed", "enhanced_speed"}),
		Cadence:   parseInt(m["cadence"]),
	}, nil
}

func parserAttr(m map[string]string, fields []string) float64 {
	for _, field := range fields {
		v, ok := m[field]
		if ok {
			return parseFloat(v)
		}
	}
	return 0
}

func parseFloat(v string) float64 {
	r, _ := strconv.ParseFloat(v, 64)
	return r
}

func parseInt(v string) int64 {
	r, _ := strconv.ParseInt(v, 10, 64)
	return r
}

func (s *Fit2Gpx) process() (*Session, error) {
	csvPath := path.Join("/tmp", fmt.Sprintf("%d.csv", time.Now().UnixMilli()))
	defer func() {
		_ = os.Remove(csvPath)
	}()
	err := s.fit2Csv(csvPath)
	if err != nil {
		s.log.Errorf("FitCSVTool.jar cmd run err %+v", err)
		return nil, err
	}

	session, err := s.parserCsv(csvPath)
	if err != nil {
		s.log.Errorf("parser csv err: %+v", err)
		return nil, err
	}
	return session, nil
}

var gpxDemo = `
<?xml version="1.0" encoding="UTF-8"?>
<gpx creator="Garmin Connect" version="1.1"
  xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/11.xsd"
  xmlns:ns3="http://www.garmin.com/xmlschemas/TrackPointExtension/v1"
  xmlns="http://www.topografix.com/GPX/1/1"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:ns2="http://www.garmin.com/xmlschemas/GpxExtensions/v3">
	<metadata>
		<author>gpxt</author>
	</metadata>
  <trk> 
  </trk> 
</gpx> 
`

func (s *Fit2Gpx) Run() error {
	if len(s.fitFile) == 0 {
		s.log.Infof("not fit file")
		return nil
	}
	session, err := s.process()
	if err != nil {
		return errors.WithStack(err)
	}
	points := session.points
	gpxData, err := gpx.ParseString(gpxDemo)
	if err != nil {
		return errors.WithStack(err)
	}
	s.log.Infof("sport type is <%s>", session.sportType)
	gpxData.Tracks[0].Type = session.sportType
	var gpxPoints []gpx.GPXPoint
	for _, p := range points {
		item := gpx.GPXPoint{
			Point: gpx.Point{
				Latitude:  p.Lat * 180 / 2147483648,
				Longitude: p.Lng * 180 / 2147483648,
				Elevation: *gpx.NewNullableFloat64(p.Altitude),
			},
			Timestamp: time.Unix(p.Ts, 0).UTC(),
		}
		exts := genExtensions(p)
		if exts != nil {
			item.Extensions = *exts
		}

		gpxPoints = append(gpxPoints, item)
	}
	gpxData.Tracks[0].Segments = []gpx.GPXTrackSegment{
		{
			Points: gpxPoints,
		},
	}

	newXml, err := gpxData.ToXml(gpx.ToXmlParams{
		Indent: true,
	})
	f, err := os.Create(s.resultGpx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	_, err = f.Write(newXml)
	if err != nil {
		return errors.WithStack(err)
	}

	s.log.Infof("fit2gpx success %s", s.resultGpx)
	return nil
}

func genExtensions(p *Record) *gpx.Extension {

	var nodes []gpx.ExtensionNode
	if p.Cadence != 0 {
		nodes = append(nodes, gpx.ExtensionNode{
			XMLName: xml.Name{
				Space: "ns3",
				Local: "cadence",
			},
			Data: fmt.Sprintf("%d", p.Cadence),
		})
	}

	if p.Altitude != 0 {
		nodes = append(nodes, gpx.ExtensionNode{
			XMLName: xml.Name{
				Space: "ns3",
				Local: "distance",
			},
			Data: fmt.Sprintf("%f", p.Distance),
		})
	}

	if p.HeartRate != 0 {
		nodes = append(nodes, gpx.ExtensionNode{
			XMLName: xml.Name{
				Space: "ns3",
				Local: "hr",
			},
			Data: fmt.Sprintf("%d", p.HeartRate),
		})
	}
	if p.Speed != 0 {
		nodes = append(nodes, gpx.ExtensionNode{
			XMLName: xml.Name{
				Space: "ns3",
				Local: "speed",
			},
			Data: fmt.Sprintf("%f", p.Speed),
		})
	}
	if len(nodes) == 0 {
		return nil
	}

	return &gpx.Extension{
		Nodes: []gpx.ExtensionNode{
			{
				XMLName: xml.Name{
					Space: "ns3",
					Local: "TrackPointExtension",
				},
				Nodes: nodes,
			}},
	}
}
