package fileutils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func PrependTo(filename string, s string) error {
	_, name := filepath.Split(filename)
	dst, err := ioutil.TempFile("", name)
	if err != nil {
		return err
	}
	defer func() {
		dst.Close()
		os.Remove(dst.Name()) // clean up
	}()

	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	return do(
		func() (err error) { _, err = io.WriteString(dst, s); return },
		func() (err error) { _, err = io.Copy(dst, src); return },
		func() error { return dst.Close() },
		func() error { return src.Close() },
		func() error { return os.Rename(dst.Name(), filename) },
	)
}

func do(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func dieIf(err error) {
	if err != nil {
		panic(err)
	}
}
