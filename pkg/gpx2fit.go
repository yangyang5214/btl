package pkg

import (
	"encoding/xml"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/tkrajina/gpxgo/gpx"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func addSpeedExt(speed float64, nodes []gpx.ExtensionNode, trackExtSpace string) []gpx.ExtensionNode {
	var (
		extNode *gpx.ExtensionNode
	)
	for i := 0; i < len(nodes); i++ {
		node := &nodes[i]
		if node.XMLName.Space == trackExtSpace {
			extNode = node
		}
	}
	speedNode := gpx.ExtensionNode{
		XMLName: xml.Name{
			Space: trackExtSpace,
			Local: "speed",
		},
		Data: fmt.Sprintf("%f", speed),
	}
	if extNode == nil {
		nodes = append(nodes, speedNode)
	} else {
		var hasSpeed bool
		for _, node := range extNode.Nodes {
			if node.XMLName.Local == "speed" {
				hasSpeed = true
				break
			}
		}
		if !hasSpeed {
			extNode.Nodes = append(extNode.Nodes, speedNode)
		}
	}
	return nodes
}

func getRatio(maxSpeedStr string, maxSpeed float64) float64 {
	var ratio = 1.0
	actualMaxSpeed := maxSpeed

	msFloat, err := strconv.ParseFloat(maxSpeedStr, 64)
	if err != nil {
		return ratio
	}
	return msFloat / actualMaxSpeed
}

func appendGpxSpeed(log *log.Helper, gpxData *gpx.GPX, maxSpeed string) ([]byte, error) {
	nsAttrs := gpxData.Attrs.GetNamespaceAttrs()

	actualMaxSpeed := gpxData.MovingData().MaxSpeed * 3.6
	ratio := getRatio(maxSpeed, actualMaxSpeed)

	log.Infof("set maxSpeed <%s km/h>, actualMaxSpeed: <%f km/h>, use ratio %f", maxSpeed, actualMaxSpeed, ratio)

	var trackExtSpace string
	for _, attr := range nsAttrs {
		if strings.Contains(attr.Value, "TrackPointExtension") {
			trackExtSpace = attr.Value
		}
	}
	for _, track := range gpxData.Tracks {
		for _, segment := range track.Segments {
			for index := range segment.Points {
				point := &segment.Points[index]
				speed := segment.Speed(index)
				point.Extensions.Nodes = addSpeedExt(speed*ratio, point.Extensions.Nodes, trackExtSpace)
			}
		}
	}
	return gpxData.ToXml(gpx.ToXmlParams{
		Indent:  true,
		Version: "1.1",
	})
}

func GenFitFile(maxSpeed, activityType string, gpx2FitCmd string, logHelper *log.Helper, gpxBytes []byte, fitFile string) error {
	//替换特殊字符
	gpxContent := strings.ReplaceAll(string(gpxBytes), "&", "&amp;")
	gpxBytes = []byte(gpxContent)

	gpxData, err := gpx.ParseBytes(gpxBytes)
	if err != nil {
		return errors.WithStack(err)
	}

	//这里统一转为 UTC 时间
	if strings.Contains(gpxData.Description, "行者") || gpxData.Description == "Export from Mi Fitness" {
		for _, track := range gpxData.Tracks {
			for _, segment := range track.Segments {
				for index := range segment.Points {
					p := &segment.Points[index] //注意这里
					p.Timestamp = p.Timestamp.Add(-time.Minute * 60 * 8)
				}
			}
		}
	}

	gpxBytes, err = appendGpxSpeed(logHelper, gpxData, maxSpeed)
	if err != nil {
		return errors.WithStack(err)
	}
	gpxFile, err := os.CreateTemp("", "*.gpx")
	if err != nil {
		logHelper.Errorf("create temp gpx err %+v", err)
		return err
	}
	logHelper.Infof("Temp gpx File %s", gpxFile.Name())
	defer func() {
		_ = os.Remove(gpxFile.Name())
	}()
	_, err = gpxFile.Write(gpxBytes)
	if err != nil {
		logHelper.Errorf("write temp gpx err %+v", err)
		return err
	}

	gpx2fitCmd := fmt.Sprintf("%s %s %s %s", gpx2FitCmd, gpxFile.Name(), fitFile, activityType)
	logHelper.Infof("run gpx2fit cmd %s", gpx2fitCmd)
	cmd := exec.Command("/bin/bash", "-c", gpx2fitCmd)

	logOutput, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := string(logOutput)
		logHelper.Errorf("run gpx2fit-java err: %s", errMsg)

		tigger := isHappy(errMsg)
		if tigger == "" {
			return errors.New("内部错误")
		} else {
			return errors.New(tigger)
		}
	}
	logHelper.Infof("gpx2fit success")
	return nil
}

func isHappy(msg string) (str string) {
	defer func() {
		str = strings.TrimSpace(str)
	}()
	msgs := strings.Split(msg, "\n")
	for _, str := range msgs {
		if strings.Contains(str, "HappyException") {
			return strings.Split(str, "HappyException:")[1]
		}
	}
	return ""
}
