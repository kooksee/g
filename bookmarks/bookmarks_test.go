package bookmarks

import (
	"bytes"
	"github.com/kooksee/g/assert"
	"io/ioutil"
	"testing"
)

func TestName(t *testing.T) {
	dt, err := ioutil.ReadFile("bookmarks.html")
	assert.MustNotError(err)

	bks, err := Parse(bytes.NewBuffer(dt))
	assert.MustNotError(err)
	assert.P(bks)
}
