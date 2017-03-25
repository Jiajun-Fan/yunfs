package main

import (
	"encoding/gob"
	"fmt"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"io"
	"time"
)

const kDefaultFsNameLength = 16
const kDefaultFsKeyLength = 16

type FileSystem struct {
	config  *Config
	oss     Oss
	cache   Oss
	entryId int
	Name    string
	entries []*Entry
}

func NewFileSystem(config *Config) *FileSystem {
	fs := &FileSystem{}
	fs.config = config
	fs.oss = MakeOss(config.Oss)
	fs.cache = MakeOss(config.Cache)
	return fs
}

func (fs *FileSystem) WriteFileEntries() {
	var fp io.WriteCloser
	var fpe Encryptor
	var enc *gob.Encoder
	for i := 0; i < len(fs.entries); i++ {
		if i%fs.config.Fs.BlockSize == 0 {
			fName := fmt.Sprintf("%s_%d", fs.config.Fs.Prefix, i/fs.config.Fs.BlockSize)
			println(fName)
			if fpe != nil {
				fpe.Close()
			}
			if fp != nil {
				fp.Close()
			}
			fp, _ = fs.oss.Create(fName)
			fpe = MakeEncryptor(fs.config.Enc, fp)
			enc = gob.NewEncoder(fpe)
		}
		entry := fs.entries[i]
		enc.Encode(entry)
	}
	if fpe != nil {
		fpe.Close()
		fp.Close()
	}
}

func (fs *FileSystem) readFileEntries() {
	index := 0
	for {
		fName := fmt.Sprintf("%s_%d", fs.config.Fs.Prefix, index)
		index++
		if err := fs.oss.Stat(fName); err != nil {
			break
		} else {
			fp, _ := fs.oss.Open(fName)
			fpd := MakeDecryptor(fs.config.Enc, fp)
			dec := gob.NewDecoder(fpd)
			for {
				entry := &Entry{}
				if err := dec.Decode(entry); err != nil {
					fpd.Close()
					fp.Close()
					break
				} else {
					if entry.Id != 0 {
						entry.fs = fs
						entry.Node = nodefs.NewDefaultNode()
						fs.entries = append(fs.entries, entry)
						if fs.entryId+1 != entry.Id {
							panic("entry ID error")
						}
						fs.entryId = entry.Id
						entry.fs = fs
					} else {
						panic("root should not in database")
					}
				}
			}
		}
	}
}

func (fs *FileSystem) String() string {
	return fs.Name
}

func (fs *FileSystem) Root() nodefs.Node {
	return fs.entries[0]
}

func (fs *FileSystem) onMount() {
	fs.readFileEntries()
	for _, entry := range fs.entries {
		if entry.Id == 0 {
			continue
		}
		pentry := fs.entries[entry.ParentId]
		pentry.Inode().NewChild(entry.Name, entry.Dir, entry)
	}
}

func (fs *FileSystem) Mount() (*fuse.Server, error) {
	root := &Entry{0, 0, true, "", "", "", fs, nodefs.NewDefaultNode(), 0}
	fs.entries = append(fs.entries, root)
	opts := &nodefs.Options{
		AttrTimeout:  time.Duration(float64(time.Second)),
		EntryTimeout: time.Duration(float64(time.Second)),
		Debug:        false,
	}
	println(fs.config.Fs.MountPoint)
	server, _, err := nodefs.MountRoot(fs.config.Fs.MountPoint, root, opts)
	return server, err
}
