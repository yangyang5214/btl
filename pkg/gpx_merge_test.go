package pkg

import (
	"testing"
)

func TestMerge(t *testing.T) {
	g := NewGpxMerge("/tmp/2")
	err := g.Run()
	if err != nil {
		t.Fatal(err)
	}
}
