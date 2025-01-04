package lazyfingertree

import (
	"fmt"
	"io"
)

// An empty finger-tree.
type emptyTree struct {
	_measurement measurement
}

func newEmptyTree(measurer measurer) fingerTree {
	return &emptyTree{measurement{measurer, measurer.Identity()}}
}

func (e *emptyTree) String() string {
	return "emptyTree{}"
}

func (e *emptyTree) Dump(w io.Writer, level int) {}

func (e *emptyTree) measurement() measurement {
	return e._measurement
}

func (e *emptyTree) AddFirst(value any) fingerTree {
	return newSingleTree(measurerFor(e), value)
}

func (e *emptyTree) AddLast(value any) fingerTree {
	return newSingleTree(measurerFor(e), value)
}

func (e *emptyTree) RemoveFirst() fingerTree {
	panic(fmt.Errorf("%w: cannot call RemoveFirst", ErrEmptyTree))
}

func (e *emptyTree) RemoveLast() fingerTree {
	panic(fmt.Errorf("%w: cannot call RemoveLast", ErrEmptyTree))
}

func (e *emptyTree) PeekFirst() any {
	panic(fmt.Errorf("%w: cannot call PeekFirst", ErrEmptyTree))
}

func (e *emptyTree) PeekLast() any {
	panic(fmt.Errorf("%w: cannot call PeekLast", ErrEmptyTree))
}

func (e *emptyTree) Concat(other fingerTree) fingerTree {
	return other
}

func (e *emptyTree) Split(pred predicate) (fingerTree, fingerTree) {
	return e, e
}

// never called but required for the interface
func (e *emptyTree) splitTree(pred predicate, initial any) (fingerTree, any, fingerTree) {
	return e, nil, e
}

func (d *emptyTree) ToSlice() []any {
	return []any{}
}

func (d *emptyTree) Each(f iterFunc) bool {
	return true
}

func (d *emptyTree) EachReverse(f iterFunc) bool {
	return true
}
