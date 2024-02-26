package amap_replaying

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tkrajina/gpxgo/gpx"
	"os"
	"strings"
)

type AmapReplaying struct {
	gpxData *gpx.GPX
	log     *log.Helper
	points  []gpx.GPXPoint
}

func NewAmapReplaying(gpxData *gpx.GPX, logger log.Logger) *AmapReplaying {
	return &AmapReplaying{
		gpxData: gpxData,
		log:     log.NewHelper(logger),
	}
}

func (r *AmapReplaying) Run() error {
	var sb strings.Builder

	sb.WriteString(r.start())
	sb.WriteString(r.addPoints())
	sb.WriteString(r.moveAnimation(1))
	sb.WriteString(r.end())

	outHtml, err := os.Create("amap_replaying.html")
	if err != nil {
		return err
	}
	defer outHtml.Close()
	_, _ = outHtml.WriteString(sb.String())
	return nil
}

func (r *AmapReplaying) addPoints() string {
	var sb strings.Builder
	sb.WriteString("var  lineArr = [\n")
	points := r.gpxData.Tracks[0].Segments[0].Points
	for i := 0; i < len(points); i += 50 {
		point := points[i]
		sb.WriteString(fmt.Sprintf("[%f,%f],\n", point.Point.Longitude, point.Point.Latitude))
	}
	r.points = points
	sb.WriteString("]")
	return sb.String()
}

func (r *AmapReplaying) start() string {
	return `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="initial-scale=1.0, user-scalable=no, width=device-width">
    <title>轨迹回放</title>
    <link rel="stylesheet" href="https://a.amap.com/jsapi_demos/static/demo-center/css/demo-center.css"/>
    <style>
        html, body, #container {
            height: 100%;
            width: 100%;
        }
    </style>
</head>
<body>
<div id="container"></div>
<script type="text/javascript" src="https://webapi.amap.com/maps?v=2.0&key=04625a30c4c1d00ab371618a37bcc59f"></script>
<script>
    window.onload = function () {
        setTimeout(function () {
            startAnimation();
        }, 5000);
    }
`
}

func (r *AmapReplaying) moveAnimation(milliSeconds int) string {
	var sb strings.Builder
	moveAnimation := fmt.Sprintf(`
    AMap.plugin('AMap.MoveAnimation', function () {
        var map = new AMap.Map("container", {
            resizeEnable: true,
            zoom: 6
        });


        var marker = new AMap.Marker({
            map: map,
            position: lineArr[0],
            icon: new AMap.Icon({
                image: "https://merge-gpx-public-1256523277.cos.ap-guangzhou.myqcloud.com/icons/biker.png",
                size: new AMap.Size(22, 22),  //图标所处区域大小
                imageSize: new AMap.Size(22,22) //图标大小
            }),
            offset: new AMap.Pixel(-13, -26),
        });

        // 绘制轨迹
        var polyline = new AMap.Polyline({
            map: map,
            path: lineArr,
            showDir: true,
            strokeColor: "#000000",  //线颜色
            strokeWeight: 8,      //线宽
        });

        var passedPolyline = new AMap.Polyline({
            map: map,
			showDir: true,
            strokeColor: "#ff0000",  //线颜色
            strokeWeight: 8,      //线宽
        });


        marker.on('moving', function (e) {
            passedPolyline.setPath(e.passedPath);
            map.setCenter(e.target.getPosition(), true)
        });

        map.setFitView();


        window.startAnimation = function startAnimation() {
            marker.moveAlong(lineArr, {
                duration: %d,
                autoRotation: true,
            });
        };
`, milliSeconds)

	sb.WriteString(moveAnimation)
	sb.WriteString(r.startEndMarker(r.points[0], r.points[len(r.points)-1]))
	sb.WriteString("\n    	});\n")
	return sb.String()
}

func (r *AmapReplaying) end() string {
	return `
</script>
</body>
</html>
`
}

func (g *AmapReplaying) startEndMarker(start, end gpx.GPXPoint) string {
	var sb strings.Builder
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

	sb.WriteString(fmt.Sprintf(`
    var startMarker = new AMap.Marker({
        position: [%f,%f],
        icon: startIcon,
        offset: new AMap.Pixel(-13, -30)
    });
`, start.Longitude, start.Latitude))

	sb.WriteString(fmt.Sprintf(`
    var endMarker = new AMap.Marker({
        position: [%f,%f],
        icon: endIcon,
        offset: new AMap.Pixel(-13, -30)
    });
`, end.Longitude, end.Latitude))
	sb.WriteString(`map.add([startMarker, endMarker]);`)
	return sb.String()
}
