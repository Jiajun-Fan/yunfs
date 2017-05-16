package main

import (
	"math/rand"
	"time"
)

func main() {
	SetDebug(DebugDebug)
	rand.Seed(time.Now().UTC().UnixNano())
	config := NewConfig()
	fs := NewFileSystem(config)
	if server, err1 := fs.Mount(); err1 != nil {
		Fatal(err1.Error())
	} else {
		server.Serve()
	}
}
