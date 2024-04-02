package gpx_amap

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/pkg/errors"
	"github.com/yangyang5214/btl/pkg"
	"github.com/yangyang5214/btl/pkg/model"
	"github.com/yangyang5214/btl/pkg/utils"
	"golang.org/x/image/colornames"
)

//https://lbs.amap.com/demo/javascript-api/example/overlayers/polyline-draw-and-edit

type GpxAmap struct {
	files []string //gpx files

	points [][]*utils.LatLng
	center utils.LatLng

	defaultColors []color.Color
	mapStyle      string
	imgPath       string
	indexHtmlPath string
	step          int
	amapKey       *model.AmapWebCode
	waitSeconds   int32

	markerStartEnd bool
	strokeWeight   int

	screenshot bool

	log    *log.Helper
	logger log.Logger
}

type TemplateAmap struct {
	Points [][]*utils.LatLng
	Center utils.LatLng
}

func NewGpxAmap(style string, logger log.Logger) *GpxAmap {
	return &GpxAmap{
		defaultColors: []color.Color{
			colornames.Red,
			colornames.Blue,
			colornames.Green,
			colornames.Black,
			colornames.Orange,
		},
		logger:         logger,
		mapStyle:       "amap://styles/" + style,
		imgPath:        "result.png",
		indexHtmlPath:  "index.html",
		step:           1,
		amapKey:        model.NewAmapWebCode(),
		markerStartEnd: true,
		strokeWeight:   8,
		screenshot:     false,
		waitSeconds:    5,
		log:            log.NewHelper(logger),
	}
}

func (g *GpxAmap) SetStrokeWeight(wight int) {
	g.strokeWeight = wight
}

func (g *GpxAmap) SetStep(step int) {
	g.step = step
}

func (g *GpxAmap) SetWaitSeconds(wait int32) {
	g.waitSeconds = wait
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

func (g *GpxAmap) Screenshot() {
	g.screenshot = true
}

func (g *GpxAmap) SetImgPath(imgPath string) {
	g.screenshot = true
	g.imgPath = imgPath
}

func (g *GpxAmap) SetColors(colors []color.Color) {
	g.defaultColors = colors
}

func (g *GpxAmap) SetMapStyle(style string) {
	g.mapStyle = style
}

func (g *GpxAmap) Run() error {
	if len(g.files) != 0 {
		for _, filename := range g.files {
			g.log.Infof("gpx file is %s", filename)
		}

		if g.amapKey == nil {
			return errors.New("amap key is required")
		}
		points, err := utils.GetPoints(g.files)
		if err != nil {
			return err
		}
		g.points = points
	}

	if len(g.points) == 0 || len(g.points[0]) == 0 {
		return errors.New("points is empty")
	}

	g.center = g.getCenter()

	g.log.Info("start gen index.html")

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

	if g.screenshot {
		g.log.Info("start screenshot")
		shot := pkg.NewScreenshot(g.imgPath, g.indexHtmlPath, g.logger)
		shot.SetWaitSeconds(g.waitSeconds)
		err = shot.Run()
		if err != nil {
			return err
		}
		g.log.Infof("screenshot success")
	}
	return nil
}

func (g *GpxAmap) randomColor(index int) string {
	relIndex := index % len(g.defaultColors)
	usedColor := g.defaultColors[relIndex]
	r := pkg.ColorToHex(usedColor)
	g.log.Infof("index %d, use color %s", index, r)
	return r
}

func (g *GpxAmap) HideStartEndPoint() {
	g.markerStartEnd = false
}

func (g *GpxAmap) startEndMarker(sb *strings.Builder) {
	sb.WriteString(`
    var startIcon = new AMap.Icon({
        size: new AMap.Size(25, 34),
        image: 'https://a.amap.com/jsapi_demos/static/demo-center/icons/dir-marker.png',
        imageSize: new AMap.Size(135, 40),
        imageOffset: new AMap.Pixel(-9, -3)
    });

    var endIcon = new AMap.Icon({
        size: new AMap.Size(25, 34),
        image: 'https://a.amap.com/jsapi_demos/static/demo-center/icons/dir-marker.png',
        imageSize: new AMap.Size(135, 40),
        imageOffset: new AMap.Pixel(-95, -3)
    });
`)

	startPoint := g.points[0][0].GCJ02String()
	sb.WriteString(fmt.Sprintf(`
    var startMarker = new AMap.Marker({
        position: new AMap.LngLat(%s),
        icon: startIcon,
        offset: new AMap.Pixel(-13, -30)
    });
`, startPoint))

	lastLine := g.points[len(g.points)-1]
	endPoint := lastLine[len(lastLine)-1].GCJ02String()

	sb.WriteString(fmt.Sprintf(`
    var endMarker = new AMap.Marker({
        position: new AMap.LngLat(%s),
        icon: endIcon,
        offset: new AMap.Pixel(-13, -30)
    });
`, endPoint))

	sb.WriteString(`map.add([startMarker, endMarker]);`)
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

			sb.WriteString(fmt.Sprintf("[%s],", point.GCJ02String()))
			//sb.WriteString(fmt.Sprintf("[%s],", point.String()))
			sb.WriteString("\n")
		}
		sb.WriteString("]")
		sb.WriteString("\n")

		polylineVar := fmt.Sprintf("polyline%d", index)
		sb.WriteString(fmt.Sprintf(`
			var %s = new AMap.Polyline({
				path: %s,
				strokeColor: "%s",
				strokeWeight: %d,
				lineJoin: 'round',
				lineCap: 'round',
			})
			%s.setMap(map)`, polylineVar, pathVar, randomColor, g.strokeWeight, polylineVar))

		sb.WriteString("\n")
	}
	if g.markerStartEnd {
		g.startEndMarker(&sb)
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

	sb.WriteString(fmt.Sprintf(`<script>
    window._AMapSecurityConfig = {
        securityJsCode: %s,
    };
</script>`, "'"+g.amapKey.Security+"'"))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf(`<script src="https://webapi.amap.com/maps?v=1.4.15&key=%s"></script>`, g.amapKey.Key))
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
