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
}

func NewAmapReplaying(gpxData *gpx.GPX) *AmapReplaying {
	return &AmapReplaying{
		gpxData: gpxData,
		log:     log.NewHelper(log.DefaultLogger),
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

	var length int
	for _, track := range r.gpxData.Tracks {
		for _, segment := range track.Segments {
			for i := 0; i < len(segment.Points); i += 30 {
				point := segment.Points[i]
				sb.WriteString(fmt.Sprintf("[%f,%f],\n", point.Point.Longitude, point.Point.Latitude))
			}
			length = length + len(segment.Points)
		}
	}
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
	return fmt.Sprintf(`
    AMap.plugin('AMap.MoveAnimation', function () {
        var map = new AMap.Map("container", {
            resizeEnable: true,
            zoom: 17
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
            strokeColor: "#3366FF",  //线颜色
            strokeWeight: 3,      //线宽
        });

        var passedPolyline = new AMap.Polyline({
            map: map,
            strokeColor: "#FF0000",  //线颜色
            strokeWeight: 3,      //线宽
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
    });
`, milliSeconds)
}
func (r *AmapReplaying) end() string {
	return `
</script>
</body>
</html>
`
}
