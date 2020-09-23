package fileutils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldPrepend(t *testing.T) {
	content := "line1\nline2"
	tmpfile, err := ioutil.TempFile("", "prepend.*.txt")
	dieIf(err)
	defer os.Remove(tmpfile.Name()) // clean up

	_, err = tmpfile.Write([]byte(content))
	dieIf(err)
	must(tmpfile.Close())

	firstLine := "foo bar\n"
	err = PrependTo(tmpfile.Name(), firstLine)
	assert.NoError(t, err)

	actual, err := ioutil.ReadFile(tmpfile.Name())
	assert.NoError(t, err)
	expected := firstLine + content
	assert.Equal(t, expected, string(actual))
}
