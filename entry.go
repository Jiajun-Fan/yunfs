package main

import (
//"fmt"
)

const (
	kNewEntry = iota
	kDeleteEntry
)

type Entry struct {
	Id       uint64
	ParentId uint64
	Dir      bool
	Name     string
	FsName   string
	Action   uint
}

func NewFile(name string, parent *Entry) *Entry {
	entry := Entry{}
	entry.Dir = false
	return &entry
}

func NewDir(name string, parent *Entry) *Entry {
	entry := Entry{}
	if parent == nil {
		entry.ParentId = 0
	}
	entry.Dir = true
	return &entry
}
