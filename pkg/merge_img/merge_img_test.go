package merge_img

import (
	"github.com/go-kratos/kratos/v2/log"
	"testing"
)

func TestMergeImg_Run(t *testing.T) {
	fg := "/Users/beer/beer/btl/pkg/merge_img/imgs/fg.png"
	bg := "/Users/beer/beer/btl/pkg/merge_img/imgs/bg.jpeg"
	mimg := NewMergeImg(fg, bg, log.DefaultLogger)
	err := mimg.Run("result.png")
	if err != nil {
		panic(err)
	}
}
