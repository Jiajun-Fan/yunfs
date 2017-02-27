package main

import (
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"strings"
)

const kPathDelimiter = "/"

type TreeNode struct {
	NodeInfo
	Parent *TreeNode
}

func NewTreeNode(entry *Entry) *TreeNode {
	node := &TreeNode{}
	node.Key = entry.Key
	node.Name = entry.Name
	node.FsName = entry.FsName
	return node
}

func (node *TreeNode) getFullName() string {
	name := node.Name

	for parent := node.Parent; parent != nil; parent = parent.Parent {
		name = fmt.Sprintf("%s%s%s", parent.Name, kPathDelimiter, name)
	}
	return name
}

func (node *TreeNode) Less(than llrb.Item) bool {
	if nthan, ok := than.(*TreeNode); !ok {
		panic("wrong type")
	} else {
		if strings.Compare(node.getFullName(), nthan.getFullName()) < 0 {
			return true
		}
	}
	return false
}
