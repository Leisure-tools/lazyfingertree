package lazyfingertree

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrFingerTree = errors.New("finger tree")

var ErrEmptyTree = fmt.Errorf("%w, empty tree", ErrFingerTree)

var ErrUnsupported = fmt.Errorf("%w, unsupported operation", ErrFingerTree)

var ErrBadMeasurer = fmt.Errorf("%w, bad measurer", ErrFingerTree)

var ErrExpectedNode = fmt.Errorf("%w, expected a node", ErrFingerTree)

type predicate func(measure any) bool

type iterFunc func(value any) bool

type HasBrief interface {
	Brief() string
}

func Brief(v any) string {
	if b, ok := v.(HasBrief); ok {
		return b.Brief()
	} else if s, ok := v.([]any); ok {
		sb := strings.Builder{}
		sb.WriteString("[")
		for i, el := range s {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(Brief(el))
		}
		sb.WriteString("]")
		return sb.String()
	}
	return strings.ReplaceAll(fmt.Sprint(v), "\n", " ")
}

type measurer interface {
	Identity() any
	// measuring a value could technically produce an error but really should not
	// make sure to validate inputs or to use a panic if you need error support
	Measure(value any) any
	Sum(a any, b any) any
}

type measurement struct {
	measurer measurer
	value    any
}

func newMeasurement(measurer measurer, value any) measurement {
	return measurement{measurer, measurer.Measure(value)}
}

func (m measurement) empty() fingerTree {
	return newEmptyTree(m.measurer)
}

// An EmptyTree, singleTree, deepTree, or delayed
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
	Each(f iterFunc) bool
	EachReverse(f iterFunc) bool
	measurement() measurement
	splitTree(predicate predicate, initial any) (fingerTree, any, fingerTree)
	fmt.Stringer
	Dump(w io.Writer, level int)
}

func isEmpty(tree fingerTree) bool {
	_, ok := force(tree).(*emptyTree)
	return ok
}

func isSingle(tree fingerTree) bool {
	_, ok := tree.(*singleTree)
	return ok
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

func takeUntil(tree fingerTree, f predicate) fingerTree {
	first, _ := tree.Split(f)
	return first
}

func dropUntil(tree fingerTree, f predicate) fingerTree {
	_, rest := tree.Split(f)
	return rest
}

// Construct a fingertree from an array.
func fromArray(measurer measurer, values []any) fingerTree {
	return prependTree(newEmptyTree(measurer), values)
}

// Prepend an array of elements to the left of a tree.
// Returns a new tree with the original one unmodified.
func prependTree[V any](tree fingerTree, values []V) fingerTree {
	for i := len(values) - 1; i >= 0; i-- {
		tree = tree.AddFirst(values[i])
	}
	return tree
}

// Append an array of elements to the right of a tree.
// Returns a new tree with the original one unmodified.
func appendTree[V any](tree fingerTree, values []V) fingerTree {
	for i := 0; i < len(values); i++ {
		tree = tree.AddLast(values[i])
	}
	return tree
}

func iterateEach(item any, f iterFunc) bool {
	if n, ok := item.(*node); ok {
		return n.Each(f)
	}
	return f(item)
}

func iterateEachReverse(item any, f iterFunc) bool {
	if n, ok := item.(*node); ok {
		return n.EachReverse(f)
	}
	return f(item)
}
