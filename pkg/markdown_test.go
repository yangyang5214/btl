package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseImages(t *testing.T) {
	m := NewMarkdown("./../README.md")
	err := m.ParseImages()
	if err != nil {
		t.Fatal(err)
	}
}

func TestParserImageUrl(t *testing.T) {
	m := Markdown{}
	r := m.parserImageUrl("![IMG_9068](https://user-images.githubusercontent.com/23392325/240136670-489412a9-9e93-461b-bc53-83e1f3fb0d47.jpeg)")
	assert.Equal(t, "https://user-images.githubusercontent.com/23392325/240136670-489412a9-9e93-461b-bc53-83e1f3fb0d47.jpeg", r)
}

func TestParserImageType(t *testing.T) {
	m := Markdown{}
	r := m.parserImageType("https://user-images.githubusercontent.com/23392325/240136670-489412a9-9e93-461b-bc53-83e1f3fb0d47.jpeg")
	assert.Equal(t, r, "jpeg")

}
