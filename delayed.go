package lazyfingertree

import (
	"fmt"
)

type fingerTreeFunc func() fingerTree

type delayed struct {
	f fingerTreeFunc
	delayedTree fingerTree
}

func newDelayed(f fingerTreeFunc) *delayed {
	tree := &delayed{f: f}
	tree.delayedTree = tree
	return tree
}

func (f *delayed) diagstr() string {
	return fmt.Sprintf("delayed{%s}", diag(f.force()))
}

func (f *delayed) force() fingerTree {
	if f.delayedTree == f {
		f.delayedTree = f.f()
	}
	return f.delayedTree
}

func (f *delayed) splitTree(predicate predicate, initial any) (fingerTree, any, fingerTree) {
	return f.force().splitTree(predicate, initial)
}

func (f *delayed) measurement() Measurement {
	return f.force().measurement()
}

func (f *delayed) AddFirst(value any) fingerTree {
	return f.force().AddFirst(value)
}

func (f *delayed) AddLast(value any) fingerTree {
	return f.force().AddLast(value)
}

func (f *delayed) RemoveFirst() fingerTree {
	return f.force().RemoveFirst()
}

func (f *delayed) RemoveLast() fingerTree {
	return f.force().RemoveLast()
}

func (f *delayed) PeekFirst() any {
	return f.force().PeekFirst()
}

func (f *delayed) PeekLast() any {
	return f.force().PeekLast()
}

func (f *delayed) Concat(other fingerTree) fingerTree {
	return f.force().Concat(other)
}

func (f *delayed) Split(predicate predicate) (fingerTree, fingerTree) {
	return f.force().Split(predicate)
}

func (f *delayed) ToSlice() []any {
	return f.force().ToSlice()
}
