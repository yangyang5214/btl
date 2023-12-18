package pkg

import (
	"golang.org/x/image/colornames"
	"testing"
)

func TestColorToHex(t *testing.T) {
	c := colornames.Yellow
	r := ColorToHex(c)
	t.Log(r)
}
