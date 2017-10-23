package main

type StackType int64

// Stack ... FILO data structure
type Stack []StackType

// Top ... return element of stack top
func (st *Stack) Top() StackType {
	if len(*st) == 0 {
		return 0
	}
	return (*st)[len(*st)-1]
}

// Push ... add element onto top of the stack
func (st *Stack) Push(x StackType) {
	*st = append(*st, x)
}

// Pop ... discard element of stack top
func (st *Stack) Pop() {
	if len(*st) != 0 {
		*st = (*st)[0 : len(*st)-1]
	}
}
