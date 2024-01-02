package model

import (
	"os"
	"path"
	"strings"
)

type AmapWebCode struct {
	Key      string
	Security string
}

func NewAmapWebCode() *AmapWebCode {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	bytes, err := os.ReadFile(path.Join(homeDir, ".amap_key"))
	if err != nil {
		return nil
	}
	data := string(bytes)
	lines := strings.Split(data, "\n")
	return &AmapWebCode{
		Key:      lines[0],
		Security: lines[1],
	}
}
