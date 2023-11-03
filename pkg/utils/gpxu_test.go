package utils

import (
	"testing"
)

func TestGetColor(t *testing.T) {
	colors := DefaultColors
	for i := 0; i < 100; i++ {
		r := GetColor(i, colors)
		t.Log(r)
	}
}
