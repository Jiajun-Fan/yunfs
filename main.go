package main

import (
	"fmt"
	"os"
)

func main() {
	SetDebug(DebugDebug)

	entries := make([]*Entry, 2)
	entries[0], _ = NewFile("ok", "cool", nil)
	entries[1], _ = NewDir("ok", "cool", nil)

	wfp, err := os.Create("entries.enc")
	if err != nil {
		panic("can't create file entries.enc")
	}

	enc, err := NewAesEncryptor([]byte("abcdefg"), wfp)
	if err != nil {
		panic("can't create encryptor")
	}

	err = WriteEntries(enc, entries)
	if err != nil {
		panic("failed to write data")
	}

	enc.Close()
	wfp.Close()

	rfp, err := os.Open("entries.enc")
	if err != nil {
		panic("can't open file entries.enc")
	}

	dec, err := NewAesDecryptor([]byte("abcdefg"), rfp)
	if err != nil {
		println(err.Error())
		panic("can't create decryptor")
	}

	defer dec.Close()
	defer rfp.Close()

	buff := make([]*Entry, 2)
	err = ReadEntries(dec, buff)
	if err != nil {
		panic(fmt.Sprintf("can't read entries %s", err.Error()))
	}

	for _, entry := range buff {
		entry.Print()
	}

}
