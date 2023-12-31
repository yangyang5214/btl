package gpx_amap

import (
	"fmt"
	"image/color"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/yangyang5214/btl/pkg"
	"github.com/yangyang5214/btl/pkg/utils"
	"golang.org/x/image/colornames"
)

var (
	// Whitesmoke 好看
	Whitesmoke = "amap://styles/whitesmoke"
	Grey       = "amap://styles/grey"

	Dark     = "amap://styles/dark"
	Light    = "amap://styles/light"
	Fresh    = "amap://styles/fresh"
	Blue     = "amap://styles/blue"
	Darkblue = "amap://styles/darkblue"
	Macaron  = "amap://styles/macaron"
)

//https://lbs.amap.com/demo/javascript-api/example/overlayers/polyline-draw-and-edit

type GpxAmap struct {
	files []string //gpx files

	points [][]*utils.LatLng
	center utils.LatLng

	defaultColors []color.Color
	mapStyle      string
	amapKey       string
	imgPath       string
	indexHtmlPath string
	step          int
}

type TemplateAmap struct {
	Points [][]*utils.LatLng
	Center utils.LatLng
}

func NewGpxAmap() *GpxAmap {
	return &GpxAmap{
		defaultColors: []color.Color{
			colornames.Red,
			colornames.Blue,
			colornames.Green,
			colornames.Black,
			colornames.Orange,
		},
		mapStyle:      Whitesmoke,
		imgPath:       "result.png",
		indexHtmlPath: "index.html",
		step:          1,
	}
}

func (g *GpxAmap) SetStep(step int) {
	g.step = step
}

func (g *GpxAmap) SetPoints(points [][]*utils.LatLng) {
	g.points = points
}

func (g *GpxAmap) SetFiles(files []string) {
	g.files = files
}

func (g *GpxAmap) SetIndexHtmlPath(indexHtmlPath string) {
	g.indexHtmlPath = indexHtmlPath
}

func (g *GpxAmap) SetImgPath(imgPath string) {
	g.imgPath = imgPath
}

func (g *GpxAmap) SetColors(colors []color.Color) {
	g.defaultColors = colors
}

func (g *GpxAmap) SetMapStyle(style string) {
	g.mapStyle = style
}

func (g *GpxAmap) loadAmapKey() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path.Join(homeDir, ".amap_key"))
	if err != nil {
		return err
	}
	g.amapKey = strings.Trim(string(data), "\n")
	return nil
}

func (g *GpxAmap) Run() error {
	err := g.loadAmapKey()
	if err != nil {
		return err
	}

	if len(g.files) != 0 {
		for _, filename := range g.files {
			log.Infof("gpx file is %s", filename)
		}

		if g.amapKey == "" {
			return errors.New("amap key is required")
		}
		points, err := utils.GetPoints(g.files)
		if err != nil {
			return err
		}
		g.points = points
	}

	if len(g.points) == 0 {
		return errors.New("points is empty")
	}

	g.center = g.getCenter()

	log.Info("start gen index.html")

	var sb strings.Builder
	sb.WriteString(g.start())
	sb.WriteString(g.drawLines())
	sb.WriteString(g.end())

	//gen index.html
	f, err := os.Create(g.indexHtmlPath)
	defer f.Close()
	_, err = f.WriteString(sb.String())
	if err != nil {
		return err
	}

	//screenshot
	log.Info("start screenshot")
	shot := pkg.NewScreenshot(g.imgPath, g.indexHtmlPath)
	return shot.Run()
}

func (g *GpxAmap) randomColor(index int) string {
	relIndex := index % len(g.defaultColors)
	r := pkg.ColorToHex(g.defaultColors[relIndex])
	//log.Infof("index %d, use color %s", index, r)
	return r
}

func (g *GpxAmap) drawLines() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`
    var map = new AMap.Map("container", {
        center: [%s],
        mapStyle: "%s"
    });
`, g.center.String(), g.mapStyle))

	sb.WriteString("\n")

	//add paths
	for index, points := range g.points {
		randomColor := g.randomColor(index)
		pathVar := fmt.Sprintf("pathIndex%d", index)
		sb.WriteString(fmt.Sprintf("var %s = [", pathVar))

		sb.WriteString("\n")

		for j := 0; j < len(points); j += g.step {
			point := points[j]
			sb.WriteString(fmt.Sprintf("[%s],", point.String()))
			sb.WriteString("\n")
		}
		sb.WriteString("]")
		sb.WriteString("\n")

		polylineVar := fmt.Sprintf("polyline%d", index)
		sb.WriteString(fmt.Sprintf(`
			var %s = new AMap.Polyline({
				path: %s,
				strokeColor: "%s",
				strokeWeight: 3,
				lineJoin: 'round',
				lineCap: 'round',
			})
			%s.setMap(map)`, polylineVar, pathVar, randomColor, polylineVar))

		sb.WriteString("\n")
	}

	sb.WriteString("map.setFitView();")
	sb.WriteString("\n")

	return sb.String()
}

func (g *GpxAmap) start() string {
	var sb strings.Builder
	sb.WriteString(`<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="initial-scale=1.0, user-scalable=no, width=device-width">
    <style>
        html,
        body,
        #container {
            width: 100%;
            height: 100%;
        }
    </style>
    <link rel="stylesheet" href="https://a.amap.com/jsapi_demos/static/demo-center/css/demo-center.css"/>
`)
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf(`<script src="https://webapi.amap.com/maps?v=1.4.15&key=%s"></script>`, g.amapKey))
	sb.WriteString("\n")

	sb.WriteString(`
    <script src="https://a.amap.com/jsapi_demos/static/demo-center/js/demoutils.js"></script>
</head>
<body>
<div id="container"></div>
<script type="text/javascript">
`)
	return sb.String()
}

func (g *GpxAmap) end() string {
	return `</script>
</body>
</html>
`
}

func (g *GpxAmap) getCenter() utils.LatLng {
	var totalLat, totalLng float64
	var count int
	for _, subs := range g.points {
		for _, point := range subs {
			totalLat = totalLat + point.Lat
			totalLng = totalLng + point.Lng
		}
		count = count + len(subs)
	}
	return utils.LatLng{
		Lat: totalLat / float64(count),
		Lng: totalLng / float64(count),
	}
}
