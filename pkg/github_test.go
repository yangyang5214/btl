package pkg

import (
	"testing"
)

func TestDownloadGist(t *testing.T) {
	github := NewGithub()
	_, err := github.DownloadGist("656e782e0c9357db5d77fea662331e1f")
	if err != nil {
		t.Fatal(err)
	}
}
