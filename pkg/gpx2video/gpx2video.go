package gpx2video

import (
	"github.com/fogleman/gg"
	"github.com/tkrajina/gpxgo/gpx"
	"math"
	"os"
)

const R = 6378137

// 解析 GPX 文件
func parseGPX(filePath string) ([]gpx.GPXPoint, error) {
	gpxFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer gpxFile.Close()

	gpxData, err := gpx.Parse(gpxFile)
	if err != nil {
		return nil, err
	}

	var points []gpx.GPXPoint
	for _, track := range gpxData.Tracks {
		for _, segment := range track.Segments {
			points = append(points, segment.Points...)
		}
	}
	return points, nil
}

// 经纬度转换为直角坐标系（墨卡托投影）
func mercatorProjection(lat, lon float64) (float64, float64) {
	x := R * lon * math.Pi / 180
	y := R * math.Log(math.Tan(math.Pi/4+lat*math.Pi/360))
	return x, y
}

// 绘制经纬度点
func plotGPX(points []gpx.GPXPoint, outputImagePath string) error {
	const width = 800
	const height = 600

	// 创建一个新的透明图像
	dc := gg.NewContext(width, height)
	dc.SetRGBA(0, 0, 0, 0) // 设置背景为完全透明
	dc.Clear()

	dc.SetRGB(1, 0, 0) // 设置轨迹点颜色为红色

	// 转换经纬度为直角坐标系
	var xPoints, yPoints []float64
	for _, point := range points {
		x, y := mercatorProjection(point.Latitude, point.Longitude)
		xPoints = append(xPoints, x)
		yPoints = append(yPoints, y)
	}

	// 计算坐标范围
	minX, maxX := xPoints[0], xPoints[0]
	minY, maxY := yPoints[0], yPoints[0]
	for i := range xPoints {
		if xPoints[i] < minX {
			minX = xPoints[i]
		}
		if xPoints[i] > maxX {
			maxX = xPoints[i]
		}
		if yPoints[i] < minY {
			minY = yPoints[i]
		}
		if yPoints[i] > maxY {
			maxY = yPoints[i]
		}
	}

	// 计算缩放比例
	scaleX := float64(width) / (maxX - minX)
	scaleY := float64(height) / (maxY - minY)
	scale := math.Min(scaleX, scaleY)

	// 绘制点
	for i := range xPoints {
		// 将坐标转换到图像空间
		x := (xPoints[i]-minX)*scale + (float64(width)-(maxX-minX)*scale)/2
		y := (float64(height) - (yPoints[i]-minY)*scale) - (float64(height)-(maxY-minY)*scale)/2
		dc.DrawPoint(x, y, 1)
		dc.Fill()
	}
	return dc.SavePNG(outputImagePath)
}
