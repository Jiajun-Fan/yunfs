package main

import (
	"encoding/gob"
	"fmt"
	"io"
)

const kDefaultFsNameLength = 16
const kDefaultFsKeyLength = 16

type FileSystem struct {
	entries []*Entry
	config  *Config
	oss     Oss
	entryId int
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

func (fs *FileSystem) FsFile(name string, fsName string, key string, parent *Entry) *Entry {
	if fsName == "" {
		fsName = string(RandStringBytes(kDefaultFsNameLength))
	}
	if key == "" {
		key = string(RandStringBytes(kDefaultFsKeyLength))
	}
	fs.entryId++
	if f, err := NewFile(fs.entryId, name, fsName, key, parent); err != nil {
		panic(err.Error())
	} else {
		fs.entries = append(fs.entries, f)
		if len(fs.entries) != fs.entryId {
			panic("append error")
		}
		return f
	}
	return nil
}

func (fs *FileSystem) FsDir(name string, fsName string, parent *Entry) *Entry {
	if fsName == "" {
		fsName = string(RandStringBytes(kDefaultFsNameLength))
	}
	fs.entryId++
	if f, err := NewDir(fs.entryId, name, fsName, parent); err != nil {
		panic(err.Error())
	} else {
		fs.entries = append(fs.entries, f)
		if len(fs.entries) != fs.entryId {
			panic("append error")
		}
		return f
	}
	return nil
}

func (fs *FileSystem) WriteFileEntries() {
	var fp io.WriteCloser
	var fpe Encryptor
	var enc *gob.Encoder
	println(len(fs.entries))
	for i := 0; i < len(fs.entries); i++ {
		if i%fs.config.Fs.BlockSize == 0 {
			fName := fmt.Sprintf("%s_%d", fs.config.Fs.Prefix, i/fs.config.Fs.BlockSize)
			if fpe != nil {
				fpe.Close()
			}
			if fp != nil {
				fp.Close()
			}
			fp, _ = fs.oss.Create(fName)
			fpe = MakeEncryptor(fs.config.Enc, fp)
			fmt.Printf("%+v\n", fs.config.Enc)
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

func (fs *FileSystem) ReadFileEntries() {
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
					entry.Print()
					fs.entries = append(fs.entries, entry)
				}
			}
		}
	}
}
