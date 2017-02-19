package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
)

const (
	kMaxEntryNumberInFile = 1024
)

type Entry struct {
	Id       uint64
	ParentId uint64
	Dir      bool
	Name     string
	FsName   string
}

func NewDir(id uint64, name string, fsName string, parent *Entry) (*Entry, error) {
	return newEntry(id, name, fsName, parent, true)
}

func NewFile(id uint64, name string, fsName string, parent *Entry) (*Entry, error) {
	return newEntry(id, name, fsName, parent, false)
}

func newEntry(id uint64, name string, fsName string, parent *Entry, dir bool) (*Entry, error) {
	entry := Entry{}
	if parent != nil && parent.Dir == false {
		return nil, errors.New("parent entry can't be file type")
	}

	entry.Id = id
	if parent == nil {
		entry.ParentId = 0
	} else {
		entry.ParentId = parent.Id
	}
	entry.Dir = dir
	entry.Name = name
	entry.FsName = fsName
	return &entry, nil
}

func ReadEntries(reader io.Reader, buff []*Entry) error {
	dec := gob.NewDecoder(reader)
	return dec.Decode(&buff)
}

func WriteEntries(writer io.Writer, buff []*Entry) error {
	enc := gob.NewEncoder(writer)
	return enc.Encode(buff)
}

func (e *Entry) Print() {
	fmt.Printf("Id: %d, ParentId: %d, Name: %s, FsName: %s, Dir: %t\n",
		e.Id, e.ParentId, e.Name, e.FsName, e.Dir)
}
