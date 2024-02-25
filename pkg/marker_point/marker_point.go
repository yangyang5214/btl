package marker_point

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tidwall/gjson"
)

type MarkerPoint struct {
	input     string
	log       *log.Helper
	showRange bool
}

type LngLat struct {
	Lng  float64
	Lat  float64
	Icon string
}

func NewMarkerPoint(input string, showRange bool) *MarkerPoint {
	return &MarkerPoint{
		input:     input,
		showRange: showRange,
		log:       log.NewHelper(log.DefaultLogger),
	}
}

func (m *MarkerPoint) Run() error {
	datas, err := os.ReadFile(m.input)
	if err != nil {
		return err
	}
	var points []*LngLat
	jsonResult := gjson.Parse(string(datas))
	for _, jsonItem := range jsonResult.Array() {
		points = append(points, &LngLat{
			Lng:  jsonItem.Get("Longitude").Float(),
			Lat:  jsonItem.Get("Latitude").Float(),
			Icon: jsonItem.Get("Icon").Str,
		})
	}
	m.log.Infof("all points is %d", len(points))

	var sb strings.Builder
	sb.WriteString(m.start())
	for i, point := range points {
		sb.WriteString(m.formatPointMarker(point, i))
		if m.showRange {
			sb.WriteString(m.addRange(point, i))
		}
	}

	sb.WriteString(m.setFitView(points))
	sb.WriteString(m.end())

	outFile, err := os.Create(fmt.Sprintf("marker_point_show_range_%v.html", m.showRange))
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, _ = outFile.WriteString(sb.String())
	return nil
}

func (m *MarkerPoint) setFitView(points []*LngLat) string {
	step := len(points) / 10
	if step < 10 {
		step = 1
	}
	var positions strings.Builder
	for i := 0; i < len(points); i += step {
		positions.WriteString(fmt.Sprintf("[%f,%f],", points[i].Lng, points[i].Lat))
	}
	return fmt.Sprintf(`
		var positions = [
			%s
		]
		var polygon = new AMap.Polygon({
			path: positions,
			map: map,
			strokeOpacity: 0,
			fillOpacity: 0,
			bubble: true
		});
		var overlaysList = map.getAllOverlays('polygon');
		map.setFitView(overlaysList);
	`, positions.String())
}

func (m *MarkerPoint) end() string {
	return `
		</script>
		</body>
		</html>
		`
}

func (m *MarkerPoint) start() string {
	return `
<!doctype html>
<html>

<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="initial-scale=1.0, user-scalable=no, width=device-width">
    <title>默认点标记</title>
    <link rel="stylesheet" href="https://a.amap.com/jsapi_demos/static/demo-center/css/demo-center.css"/>
    <style>
        html,
        body,
        #container {
            height: 100%;
            width: 100%;
        }
    </style>
</head>

<body>
<div id="container"></div>
<script type="text/javascript"
        src="https://webapi.amap.com/maps?v=1.4.15&key=d681abe16c6b76a8da6d2762ad881134"></script>
<script type="text/javascript">
    var marker, map = new AMap.Map("container", {
        resizeEnable: true,
        zoom: 13
    });
`
}

func (m *MarkerPoint) addRange(point *LngLat, index int) string {
	return fmt.Sprintf(`
var circle%d = new AMap.Circle({
        center: [%f,%f],
        radius: 500,
        borderWeight: 3,
        strokeColor: "#FF33FF", 
        strokeOpacity: 0,
        strokeWeight: 6,
        fillOpacity: 0.4,
        strokeStyle: 'dashed',
        strokeDasharray: [10, 10], 
        fillColor: '#1791fc',
        zIndex: 50,
    })
    circle%d.setMap(map)
`, index, point.Lng, point.Lat, index)
}

func (m *MarkerPoint) formatPointMarker(point *LngLat, index int) string {
	icon := point.Icon
	if icon == "" {
		icon = "https://a.amap.com/jsapi_demos/static/demo-center/icons/poi-marker-default.png"
	}
	return fmt.Sprintf(`
    marker%d = new AMap.Marker({
        position: [%f, %f],
		icon: new AMap.Icon({
			 image: "%s",
			 size: new AMap.Size(22, 28),  //图标所处区域大小
			 imageSize: new AMap.Size(22,28) //图标大小
		})	
    });
    marker%d.setMap(map);
`, index, point.Lng, point.Lat, icon, index)
}
