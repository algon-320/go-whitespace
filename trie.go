package main

import (
	"fmt"
	"strings"
)

type trieNode struct {
	child map[rune]*trieNode
	imp   string
	ins   Instruction
}

func NewTrie(dict map[string]Instruction) *trieNode {
	root := new(trieNode)
	root.child = make(map[rune]*trieNode)
	root.imp = ""
	root.ins = nil

	for imp, ins := range dict {
		cur := root
		for _, c := range imp {
			next, exist := cur.child[c]
			if exist {
				cur = next
			} else {
				cur.child[c] = new(trieNode)
				cur.child[c].child = make(map[rune]*trieNode)
				cur.child[c].ins = nil
				cur.child[c].imp = cur.imp + string(c)
				cur = cur.child[c]
			}
		}
		cur.ins = ins
	}
	return root
}

func (n *trieNode) Print(depth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Print(indent, n.imp, ", ")
	if n.ins != nil {
		fmt.Print(n.ins.Name())
	}
	fmt.Print("\n")
	for _, c := range n.child {
		c.Print(depth + 1)
	}
}

func (n *trieNode) Find(c rune) *trieNode {
	next, exist := n.child[c]
	if exist {
		return next
	}
	return nil
}

var trieRoot *trieNode

func initTrie() {
	dict := make(map[string]Instruction)
	for _, v := range InstructionList {
		dict[v.Imp()] = v
	}
	trieRoot = NewTrie(dict)
	// triRoot.Print(0)
}
