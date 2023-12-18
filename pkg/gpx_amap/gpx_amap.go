package gpx_amap

import (
	"fmt"
	"github.com/yangyang5214/btl/pkg"
	"github.com/yangyang5214/btl/pkg/utils"
	"golang.org/x/image/colornames"
	"image/color"
	"os"
	"strings"
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
	allColorNames []string // color_names
	mapStyle      string
}

type TemplateAmap struct {
	Points [][]*utils.LatLng
	Center utils.LatLng
}

func NewGpxAmap(files []string) *GpxAmap {
	return &GpxAmap{
		files:         files,
		allColorNames: colornames.Names,
		mapStyle:      Whitesmoke,
	}
}

func (g *GpxAmap) SetColors(colors []color.Color) {
	g.defaultColors = colors
}

func (g *GpxAmap) SetMapStyle(style string) {
	g.mapStyle = style
}

func (g *GpxAmap) Run() error {
	points, err := utils.GetPoints(g.files)
	if err != nil {
		return err
	}
	g.points = points
	g.center = g.getCenter()

	var sb strings.Builder
	sb.WriteString(g.start())
	sb.WriteString(g.drawLines())
	sb.WriteString(g.end())

	f, err := os.Create("index.html")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(sb.String())
	if err != nil {
		return err
	}
	return nil
}

func (g *GpxAmap) randomColor(index int) string {
	if len(g.defaultColors) == 0 {
		for _, v := range colornames.Map {
			g.defaultColors = append(g.defaultColors, v)
		}
	}
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

		for _, point := range points {
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
	return `<!doctype html>
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
    <script src="https://webapi.amap.com/maps?v=1.4.15&key=04625a30c4c1d00ab371618a37bcc59f"></script>
    <script src="https://a.amap.com/jsapi_demos/static/demo-center/js/demoutils.js"></script>
</head>
<body>
<div id="container"></div>
<script type="text/javascript">
`
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
