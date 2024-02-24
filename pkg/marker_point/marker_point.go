package marker_point

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tidwall/gjson"
)

type MarkerPoint struct {
	input string
	log   *log.Helper
}

type LngLat struct {
	Lng float64
	Lat float64
}

func NewMarkerPoint(input string) *MarkerPoint {
	return &MarkerPoint{
		input: input,
		log:   log.NewHelper(log.DefaultLogger),
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
			Lng: jsonItem.Get("Longitude").Float(),
			Lat: jsonItem.Get("Latitude").Float(),
		})
	}
	m.log.Infof("all points is %d", len(points))

	var sb strings.Builder
	sb.WriteString(m.start())
	for i, point := range points {
		mp := m.formatPointMarker(point, i)
		sb.WriteString(mp)
	}
	sb.WriteString(m.end())

	outFile, err := os.Create("marker_point.html")
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, _ = outFile.WriteString(sb.String())
	return nil
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

func (m *MarkerPoint) formatPointMarker(point *LngLat, index int) string {
	return fmt.Sprintf(`
    marker%d = new AMap.Marker({
        position: [%f, %f],
    });
    marker%d.setMap(map);
`, index, point.Lng, point.Lat, index)
}
