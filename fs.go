package main

type FileSystem struct {
	entries []Entry
	//root     *Entry
	curIndex uint64
	maxIndex uint64
}

func NewFileSystem() *FileSystem {
	fs := &FileSystem{}
	fs.maxIndex = 256 * 1024
	fs.entries = make([]Entry, 0, fs.maxIndex)
	return fs
}

func (fs *FileSystem) AddFile() {
}
