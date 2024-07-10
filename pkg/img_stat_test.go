package pkg

import (
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestStat(t *testing.T) {
	stat := NewImgStat("/tmp/1,png", log.DefaultLogger)

	var infos []*StatInfo

	infos = append(infos, &StatInfo{
		Label: "距离",
		Value: "100 KM",
	})

	infos = append(infos, &StatInfo{
		Label: "时间",
		Value: "100h 20m",
	})

	err := stat.Run(infos)
	if err != nil {
		panic(err)
	}
}
