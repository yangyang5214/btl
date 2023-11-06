package pkg

import (
	"testing"

	"github.com/golang/geo/s2"
)

func Test1(t *testing.T) {
	//point := []float64{108.989234, 34.18037} //西安
	//point := []float64{121.455708, 31.249574} //上海
	//point := []float64{120.207692, 30.240262} //杭州
	point := []float64{109.572927, 19.466784} //儋州
	center := s2.LatLngFromDegrees(point[1], point[0])
	for i := 8; i <= 13; i++ {
		c := NewCarto("carto-dark", i, center, 200)
		err := c.Run()
		if err != nil {
			t.Fatal(err)
		}
	}
}
