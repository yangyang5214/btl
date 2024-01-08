package pkg

import (
	"fmt"
	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
	"github.com/tkrajina/gpxgo/gpx"
)

type GpxGrid struct {
}

func NewGpxGrid() *GpxGrid {
	return &GpxGrid{}
}

func (g *GpxGrid) Run() error {
	// 解析GPX数据
	gpxData, err := gpx.ParseFile("/Users/beer/beer/rides/hai_nan.gpx")
	if err != nil {
		log.Fatal(err)
	}

	// 创建绘图上下文
	const width = 800
	const height = 600
	dc := gg.NewContext(width, height)

	// 将GPX轨迹添加到图形中
	dc.SetRGB(1, 0, 0) // 设置颜色为红色
	for _, track := range gpxData.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				x, y := point.Longitude, point.Latitude
				x, y = convertToXY(x, y, width, height)
				dc.DrawPoint(x, y, 2)
				dc.Stroke()
			}
		}
	}

	// 保存图形到文件
	if err := dc.SavePNG("output.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("图形已保存为output.png")
	return nil
}

// convertToXY将经纬度转换为图形坐标
func convertToXY(lon, lat, width, height float64) (float64, float64) {
	x := (lon + 180) * (width / 360)
	y := (90 - lat) * (height / 180)
	return x, y
}
