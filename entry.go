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

type Entry struct {
	Id       int
	ParentId int
	Dir      bool
	Key      string
	Name     string
	FsName   string
	fs       *FileSystem
	nodefs.Node
	size uint64
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
	entry := Entry{id, 0, dir, key, name, fsName, nil, nil, 0}

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
	if e.size != 0 {
		out.Size = e.fs.oss.Size(e.FsName)
		e.size = out.Size
	} else {
		out.Size = e.size
	}
}

func (e *Entry) Data() (data []byte) {
	fp, err := e.fs.oss.Open(e.FsName)
	if err != nil {
		return nil
	}
	fpd := MakeDecryptor(e.fs.config.Enc, fp)
	data, _ = ioutil.ReadAll(fpd)

	out := fuse.Attr{}
	e.Stat(&out)
	size := int(out.Size)
	if len(data) < size {
		data = append(data, make([]byte, size-len(data))...)
	}

	fpd.Close()
	fp.Close()
	return
}

func (n *Entry) OnMount(c *nodefs.FileSystemConnector) {
	n.fs.onMount()
}

func (n *Entry) OpenDir(context *fuse.Context) (stream []fuse.DirEntry, code fuse.Status) {
	children := n.Inode().Children()
	stream = make([]fuse.DirEntry, 0, len(children))
	for k, v := range children {
		mode := fuse.S_IFREG | 0666
		if v.IsDir() {
			mode = fuse.S_IFDIR | 0777
		}
		stream = append(stream, fuse.DirEntry{
			Name: k,
			Mode: uint32(mode),
		})
	}
	return stream, fuse.OK
}

func (n *Entry) Open(flags uint32, context *fuse.Context) (fuseFile nodefs.File, code fuse.Status) {
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}

	return nodefs.NewDataFile(n.Data()), fuse.OK
}

func (n *Entry) Deletable() bool {
	return false
}

func (n *Entry) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) fuse.Status {
	if n.Inode().IsDir() {
		out.Mode = fuse.S_IFDIR | 0777
		return fuse.OK
	}
	n.Stat(out)
	return fuse.OK
}
