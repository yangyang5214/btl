package pkg

import "testing"

func TestGetRatio(t *testing.T) {
	t.Run("", func(t *testing.T) {
		r := getRatio("10", 11) //变小
		t.Log(r)
	})

	t.Run("", func(t *testing.T) {
		r := getRatio("20", 10)
		t.Log(r)
	})

}
