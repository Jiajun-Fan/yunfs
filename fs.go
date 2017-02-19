package main

import (
	"fmt"
)

type FileSystem struct {
	entries []Entry
	config  *Config
	oss     Oss
	//root     *Entry
	//curIndex uint64
	//maxIndex uint64
}

func NewFileSystem(config *Config) *FileSystem {
	fs := &FileSystem{}
	fs.config = config
	fs.oss = MakeOss(config.Oss)
	//fs.maxIndex = config.Fs.CacheSize
	//fs.entries = make([]Entry, 0, fs.maxIndex)
	return fs
}

func (fs *FileSystem) readFileEntry() {
}

func (fs *FileSystem) ReadFileEntries() {
	for {
		index := 0
		fName := fmt.Sprintf("%s_%d", fs.config.Fs.Prefix, index)
		if err := fs.oss.Stat(fName); err != nil {
			break
		} else {
		}
	}
}
