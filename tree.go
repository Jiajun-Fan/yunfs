package main

import (
	"fmt"
	"strings"
)

const kPathDelimiter = "/"

type TreeNode struct {
	NodeInfo
	Parent * TreeNode
}

func (node *TreeNode) getFullName() string {
	name := node.Name
	
	for parent := node.Parent; parent != nil; parent = parent.Parent {
		name = fmt.Sprintf("%s%s%s", parent.Name, kPathDelimiter, name)
	}
	return name
}

func lessThan(a, b interface{}) bool {
	na, oka := a.(*TreeNode)
	nb, okb := b.(*TreeNode)
	if !oka || !okb {
		panic("type error")
	}
	if strings.Compare(na.getFullName(), nb.getFullName()) < 0 {
		return true
	}
	return false
}
