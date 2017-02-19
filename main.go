package main

import (
//"fmt"
//"os"
)

func main() {
	SetDebug(DebugDebug)
	if config, err := NewConfig(); err != nil {
		Fatal("no configuration file found, a template is generated at '~/.yunfs/yunfs.json'\n")
	} else {
		fs := NewFileSystem(config)
		fs.ReadFileEntry()
	}
}
