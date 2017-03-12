package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"io"
	"io/ioutil"
)

const (
	kMaxEntryNumberInFile = 1024
)

type NodeInfo struct {
	Id       int
	ParentId int
	Dir      bool
	Key      string
	Name     string
	FsName   string
}

type Entry struct {
	NodeInfo
	memNode
}

func NewDir(id int, name string, parent *Entry) (*Entry, error) {
	return newEntry(id, name, "", "", parent, true)
}

func NewFile(id int, name string, fsName string, key string, parent *Entry) (*Entry, error) {
	return newEntry(id, name, fsName, key, parent, false)
}

func newEntry(id int, name string, fsName string, key string, parent *Entry, dir bool) (*Entry, error) {
	if parent != nil && parent.Dir == false {
		return nil, errors.New("parent entry can't be file type")
	}
	entry := Entry{NodeInfo{id, 0, dir, key, name, fsName}, memNode{nodefs.NewDefaultNode(), nil, nil}}

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

func (e *Entry) Stat(out *fuse.Attr) {
	out.Mode = fuse.S_IFREG | 0444
}

func (e *Entry) Data() (data []byte) {
	fp, err := e.fs.oss.Open(e.FsName)
	if err != nil {
		return nil
	}
	fpd := MakeDecryptor(e.fs.config.Enc, fp)
	data, _ = ioutil.ReadAll(fpd)
	fpd.Close()
	fp.Close()
	return
}
