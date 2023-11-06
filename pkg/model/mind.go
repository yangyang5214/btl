package model

import (
	"fmt"
	"strings"
)

type Mind struct {
	Key   string
	Name  string
	Count int
	Level int
}

func (m *Mind) String() string {
	s := strings.Builder{}
	for i := 0; i < m.Level+1; i++ {
		s.WriteString("#")
	}
	s.WriteString(" ")

	s.WriteString(fmt.Sprintf("%s (%d)", m.Name, m.Count))

	s.WriteString("\n")
	s.WriteString("\n")
	return s.String()
}

type Bounds struct {
	X []int
	Y []int
}
