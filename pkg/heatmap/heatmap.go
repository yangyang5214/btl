package heatmap

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tidwall/gjson"
	"os"
	"strings"
)

type HeatMap struct {
	userId   string
	mergeGpx *MergeGpx
	log      *log.Helper
}

func NewHeatMap(userId string) *HeatMap {
	return &HeatMap{
		userId:   userId,
		mergeGpx: NewMergeGpx(),
		log:      log.NewHelper(log.DefaultLogger),
	}
}

func (h *HeatMap) start() string {
	return `
<!doctype html>
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
    <script>
        window._AMapSecurityConfig = {
            securityJsCode: 'a384e97096fafd0e65dd8db52a8de562',
        };
    </script>
    <script src="https://webapi.amap.com/maps?v=1.4.15&key=04625a30c4c1d00ab371618a37bcc59f"></script>
    <script src="https://webapi.amap.com/ui/1.1/main.js?v=1.1.2"></script>
</head>

<body>
<div id="container"></div>
<script>
`
}

func (h *HeatMap) lines() string {
	var sb strings.Builder
	sb.WriteString(`
 var map = new AMap.Map('container', {
        zoom: 8
    });

    //加载PathSimplifier，loadUI的路径参数为模块名中 'ui/' 之后的部分
    AMapUI.load(['ui/misc/PathSimplifier'], function (PathSimplifier) {
        if (!PathSimplifier.supportCanvas) {
            alert('当前环境不支持 Canvas！');
            return;
        }
        //启动页面
        initPage(PathSimplifier);
    });

    function initPage(PathSimplifier) {
        //创建组件实例
        var pathSimplifierIns = new PathSimplifier({
            zIndex: 100,
            map: map, //所属的地图实例
            getPath: function (pathData, pathIndex) {
                //返回轨迹数据中的节点坐标信息，[AMap.LngLat, AMap.LngLat...] 或者 [[lng|number,lat|number],...]
                return pathData.path;
            },
            getHoverTitle: function (pathData, pathIndex, pointIndex) {
                //返回鼠标悬停时显示的信息
                if (pointIndex >= 0) {
                    //鼠标悬停在某个轨迹节点上
                    return pathData.name + '，点:' + pointIndex + '/' + pathData.path.length;
                }
                //鼠标悬停在节点之间的连线上
                return pathData.name + '，点数量' + pathData.path.length;
            },
            renderOptions: {
                pathLineStyle: {
                    strokeStyle: 'red',
                    lineWidth: 1,
                    dirArrowStyle: true
                }
            }
        });
`)

	sb.WriteString(h.setLines())

	sb.WriteString(`
    }
`)

	return sb.String()
}

func (h *HeatMap) setLines() string {
	var sb strings.Builder
	sb.WriteString(`
pathSimplifierIns.setData([
`)
	ids, err := h.mergeGpx.GetActivityIds()
	if err != nil {
		panic(err)
	}

	ids = ids[:5]

	h.log.Infof("all line size %d", len(ids))
	for _, id := range ids {
		h.log.Infof("process activity_id %s", id)
		line, err := h.formatLine(id)
		if err != nil {
			panic(err)
		}
		sb.WriteString(line)
	}

	sb.WriteString("        ]);\n")
	return sb.String()
}

func (h *HeatMap) formatLine(activityId string) (string, error) {
	resp, err := h.mergeGpx.DownloadActivity(activityId)
	if err != nil {
		return "", err
	}
	result := gjson.ParseBytes(resp)

	name := result.Get("name").String()
	var path []string
	points := result.Get("path").Array()
	for _, p := range points {
		point := p.Get("point").Array()
		path = append(path, fmt.Sprintf("[%f,%f]", point[1].Float(), point[0].Float()))
	}

	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("		name: '%s',\n", name))
	sb.WriteString(fmt.Sprintf(" path: [%v]\n", strings.Join(path, ",")))
	sb.WriteString("},\n")
	return sb.String(), nil
}

func (h *HeatMap) end() string {
	return `
</script>
</body>
</html>
`
}

func (h *HeatMap) Run() error {
	var sb strings.Builder
	sb.WriteString(h.start())
	sb.WriteString(h.lines())
	sb.WriteString(h.end())

	f, err := os.Create("heatmap.html")
	if err != nil {
		return err
	}
	defer f.Close()
	_, _ = f.WriteString(sb.String())
	return nil
}
