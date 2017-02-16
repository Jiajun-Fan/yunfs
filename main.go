package main

import (
	"bytes"
	"fmt"
)

func main() {
	SetDebug(DebugDebug)

	var network bytes.Buffer
	enc, err := NewAesEncryptor([]byte(`abcdefg`), &network)
	if err != nil {
		Fatal(err.Error())
	}

	enc.Write([]byte(`this is cool  dsafasdas dfas`))

	dec, err := NewAesDecryptor([]byte(`abcdefg`), &network)
	if err != nil {
		Fatal(err.Error())
	}

	buff := make([]byte, 4096)
	n, err := dec.Read(buff)
	if err != nil {
		println(n)
		Fatal(err.Error())
	}
	fmt.Printf("%s\n", string(buff))
	println(n)
}
