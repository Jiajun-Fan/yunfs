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
	Id       int
	ParentId int
	Dir      bool
	Key      string
	Name     string
	FsName   string
}

func NewDir(id int, name string, fsName string, parent *Entry) (*Entry, error) {
	return newEntry(id, name, fsName, "", parent, true)
}

func NewFile(id int, name string, fsName string, key string, parent *Entry) (*Entry, error) {
	return newEntry(id, name, fsName, key, parent, false)
}

func newEntry(id int, name string, fsName string, key string, parent *Entry, dir bool) (*Entry, error) {
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
	entry.Key = key
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
	if e.Dir {
		fmt.Printf("[Dir]: Id: %d, ParentId: %d, Name: %s, FsName: %s\n",
			e.Id, e.ParentId, e.Name, e.FsName)
	} else {
		fmt.Printf("Id: %d, ParentId: %d, Name: %s, FsName: %s, Key %s\n",
			e.Id, e.ParentId, e.Name, e.FsName, e.Key)
	}
}
