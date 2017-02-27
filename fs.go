package main

import (
	"encoding/gob"
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	//"io"
)

const kDefaultFsNameLength = 16
const kDefaultFsKeyLength = 16

type FileSystem struct {
	config  *Config
	oss     Oss
	entryId int
	tree    *llrb.LLRB
}

func NewFileSystem(config *Config) *FileSystem {
	fs := &FileSystem{}
	fs.config = config
	fs.oss = MakeOss(config.Oss)
	fs.tree = llrb.New()
	return fs
}

/*func (fs *FileSystem) FsFile(name string, fsName string, key string, parent *Entry) *Entry {
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
}*/

func (fs *FileSystem) readFileEntries() []*Entry {
	index := 0
	entries := make([]*Entry, 0, 0)
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
					//entry.Print()
					entries = append(entries, entry)
				}
			}
		}
	}
	return entries
}

func (fs *FileSystem) BuildFileTree() {
	entries := fs.readFileEntries()
	nodes := make([]*TreeNode, len(entries)+1)

	root := &TreeNode{}
	nodes[0] = root

	for _, entry := range entries {
		node := NewTreeNode(entry)
		node.Parent = nodes[entry.ParentId]
		nodes[entry.Id] = node
		fs.tree.ReplaceOrInsert(node)
	}

	iter := func(i llrb.Item) bool {
		n, _ := i.(*TreeNode)
		max, _ := fs.tree.Max().(*TreeNode)
		println(n.getFullName())
		if n.getFullName() == max.getFullName() {
			return false
		}
		return true
	}
	fs.tree.AscendGreaterOrEqual(fs.tree.Min(), iter)
}
