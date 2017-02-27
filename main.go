package main

import (
	//"fmt"
	//"os"
	"math/rand"
	"time"
)

func main() {
	SetDebug(DebugDebug)
	rand.Seed(time.Now().UTC().UnixNano())
	if config, err := NewConfig(); err != nil {
		Fatal("no configuration file found, a template is generated at '~/.yunfs/yunfs.json'\n")
	} else {
		fs := NewFileSystem(config)
		fs.BuildFileTree()
	}
}
