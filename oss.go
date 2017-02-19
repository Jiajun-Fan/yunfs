package main

import (
	"fmt"
	"os"
)

type Oss interface {
	Stat(fname string) error
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

func MakeOss(config OssConfig) Oss {
	if config.Type == "local" {
		return NewLocalOss(config.Base)
	} else {
		Fatal(fmt.Sprintf("oss type '%s' is not implemented", config.Type))
		return nil
	}
}
