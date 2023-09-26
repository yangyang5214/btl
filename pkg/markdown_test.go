package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	r := m.parserMarkdownImgUrl("![IMG_9068](https://user-images.githubusercontent.com/23392325/240136670-489412a9-9e93-461b-bc53-83e1f3fb0d47.jpeg)")
	assert.Equal(t, "https://user-images.githubusercontent.com/23392325/240136670-489412a9-9e93-461b-bc53-83e1f3fb0d47.jpeg", r)
}

func TestParserImageType(t *testing.T) {
	m := Markdown{}
	r := m.parserImageType("https://user-images.githubusercontent.com/23392325/240136670-489412a9-9e93-461b-bc53-83e1f3fb0d47.jpeg")
	assert.Equal(t, r, "jpeg")

}

func TestParserSrcImgUrl(t *testing.T) {
	s := `<img width="994" alt="image" src="https://user-images.githubusercontent.com/23392325/270086663-88b3f1fe-6262-4243-bf95-51d2439c627f.png">`
	r := parserSrcImgUrl(s)
	t.Logf(r)
}
