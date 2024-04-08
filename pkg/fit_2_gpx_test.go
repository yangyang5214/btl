package pkg

import "testing"

func Test1(t *testing.T) {
	r, err := ParserPyResult("/tmp/1")
	if err != nil {
		panic(err)
	}
	t.Log(r)
	t.Log(len(r))
}
