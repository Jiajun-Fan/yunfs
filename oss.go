package main

import (
	"fmt"
	"io"
	"os"
)

type Oss interface {
	Stat(fname string) error
	Create(name string) (io.WriteCloser, error)
	Open(name string) (io.ReadCloser, error)
}

type LocalOss struct {
	Base string
}

func NewLocalOss(base string) Oss {
	oss := &LocalOss{}
	oss.Base = base
	return oss
}

func (l *LocalOss) Stat(fname string) error {
	fn := fmt.Sprintf("%s/%s", l.Base, fname)
	_, err := os.Stat(fn)
	return err
}

func (l *LocalOss) Create(name string) (io.WriteCloser, error) {
	fn := fmt.Sprintf("%s/%s", l.Base, name)
	return os.Create(fn)
}

func (l *LocalOss) Open(name string) (io.ReadCloser, error) {
	fn := fmt.Sprintf("%s/%s", l.Base, name)
	return os.Open(fn)
}

func MakeOss(config OssConfig) Oss {
	if config.Type == "local" {
		return NewLocalOss(config.Base)
	} else {
		Fatal(fmt.Sprintf("oss type '%s' is not implemented", config.Type))
		return nil
	}
}
