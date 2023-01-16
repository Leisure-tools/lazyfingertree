package lazyfingertree

import (
	"errors"
	"fmt"
)

var ErrFingerTree = errors.New("finger tree")

var ErrEmptyTree = fmt.Errorf("%w, empty tree", ErrFingerTree)

var ErrUnsupported = fmt.Errorf("%w, unsupported operation", ErrFingerTree)

var ErrBadMeasurer = fmt.Errorf("%w, bad measurer", ErrFingerTree)

var ErrExpectedNode = fmt.Errorf("%w, expected a node", ErrFingerTree)

type predicate func(measure any) bool

type iterFunc func(value any) bool

type measurer interface {
	Identity() any
	// measuring a value could technically produce an error but really should not
	// make sure to validate inputs or to use a panic if you need error support
	Measure(value any) any
	Sum(a any, b any) any
}

type diagable interface {
	diagstr() string
}

func diag(v any) string {
	if d, ok := v.(diagable); !ok {
		return fmt.Sprintf("%v", v)
	} else {
		return d.diagstr()
	}
}

type PluggableMeasurer struct {
	identity any
	measure func(value any) any
	sum func(a any, b any) any
}

func (m PluggableMeasurer) Identity() any {
	return m.identity
}

func (m PluggableMeasurer) Measure(v any) any {
	return m.measure(v)
}

func (m PluggableMeasurer) Sum(a any, b any) any {
	return m.sum(a, b)
}

type Measurement struct {
	measurer measurer
	value any
}

func newMeasurement(measurer measurer, m any) Measurement {
	return Measurement{measurer, m}
}

func measurement(measurer measurer, v any) Measurement {
	return newMeasurement(measurer, measurer.Measure(v))
}

func (m Measurement) empty() fingerTree {
	return newEmptyTree(m.measurer)
}

// An EmptyTree, singleTree, or deepTree
type fingerTree interface {
	AddFirst(value any) fingerTree
	AddLast(value any) fingerTree
	RemoveFirst() fingerTree
	RemoveLast() fingerTree
	PeekFirst() any
	PeekLast() any
	Concat(other fingerTree) fingerTree
	Split(predicate predicate) (fingerTree, fingerTree)
	ToSlice() []any
	measurement() Measurement
	splitTree(predicate predicate, initial any) (fingerTree, any, fingerTree)
}

func IsEmpty(tree fingerTree) bool {
	_, ok := force(tree).(*emptyTree)
	return ok
}

func isSingle(tree fingerTree) bool {
	_, ok := tree.(*singleTree)
	return ok
}

func Measure(tree fingerTree) any {
	return tree.measurement().value
}

func measurerFor(tree fingerTree) measurer {
	return tree.measurement().measurer
}

func force(tree fingerTree) fingerTree {
	t, ok := tree.(*delayed)
	if !ok {
		return tree
	}
	return t.force()
}

func empty(tree fingerTree) fingerTree {
	return newEmptyTree(tree.measurement().measurer)
}

func TakeUntil(tree fingerTree, f predicate) (fingerTree) {
	first, _ := tree.Split(f)
	return first
}

func DropUntil(tree fingerTree, f predicate) (fingerTree) {
	_, rest := tree.Split(f)
	return rest
}

func Each(tree fingerTree, f iterFunc) error {
	for (!IsEmpty(tree)) {
		if !f(tree.PeekFirst()) {
			break
		}
		tree = tree.RemoveFirst()
	}
	return nil
}

func EachReverse(tree fingerTree, f iterFunc) error {
	for (!IsEmpty(tree)) {
		if !f(tree.PeekLast()) {
			break
		}
		tree = tree.RemoveLast()
	}
	return nil
}

// Construct a fingertree from an array.
func fromArray(measurer measurer, values []any) fingerTree {
	return prependTree(newEmptyTree(measurer), values)
}

// Prepend an array of elements to the left of a tree.
// Returns a new tree with the original one unmodified.
func prependTree(tree fingerTree, values []any) fingerTree {
	for i := len(values) - 1; i >= 0; i-- {
		tree = tree.AddFirst(values[i])
	}
	return tree
}

// Append an array of elements to the right of a tree.
// Returns a new tree with the original one unmodified.
func appendTree(tree fingerTree, values []any) fingerTree {
	for i := 0; i < len(values); i++ {
		tree = tree.AddLast(values[i])
	}
	return tree
}

func concat(slices ...[]any) []any {
	size := 0
	for _, slice := range slices {
		size += len(slice)
	}
	result := make([]any, 0, size)
	for slice := range slices {
		result = append(result, slice)
	}
	return result
}
