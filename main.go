package main

import (
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
		//fs.readFileEntries()
		//fs.WriteFileEntries()
		if server, err1 := fs.Mount(); err1 != nil {
			Fatal(err1.Error())
		} else {
			server.Serve()
		}
	}
}
