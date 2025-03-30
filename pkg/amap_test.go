package pkg

import (
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestGeoParser(t *testing.T) {
	amap := NewAmap(log.NewHelper(log.DefaultLogger))

	r, err := amap.GetLocationByAddress("上海市浦东新区陆家嘴街道")
	if err != nil {
		panic(err)
	}
	t.Log(r)
}
