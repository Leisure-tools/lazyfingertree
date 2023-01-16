package lazyfingertree

import "fmt"

// An empty finger-tree.
type emptyTree struct {
	_measurement Measurement
}

func newEmptyTree(measurer measurer) fingerTree {
	return &emptyTree{newMeasurement(measurer, measurer.Identity())}
}

func (e *emptyTree) diagstr() string {
	return "emptyTree{}"
}

func (e *emptyTree) measurement() Measurement {
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
