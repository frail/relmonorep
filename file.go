package main

import (
	"io"
	"io/ioutil"
	"os"
)

type file struct {
	name    string
	writer  io.WriteCloser
	content []byte
}

func newFile(name string) (*file, error) {
	f := &file{name: name, writer: os.Stdout}
	var err error
	if name == "-" {
		return f, nil
	}

	f.content, err = ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	f.writer, err = os.OpenFile(name, os.O_RDWR|os.O_TRUNC, 0666)
	return f, err

}

func (f *file) Write(p []byte) (n int, err error) {
	return f.writer.Write(p)
}

func (f *file) Close() error {
	if f.name == "-" || f.writer == nil {
		return nil
	}

	if f.content != nil {
		f.writer.Write(f.content)
	}
	return f.writer.Close()
}
