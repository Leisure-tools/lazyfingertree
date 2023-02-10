package lazyfingertree

import (
	"fmt"
	"strings"
)

// A node is a measured container of either 2 or 3 sub-finger-trees.
type node struct {
	_measurement measurement
	children     []any
}

func asNode(v any) *node {
	if n, ok := v.(*node); !ok {
		panic(ErrExpectedNode)
	} else {
		return n
	}
}

func newNode(measurer measurer, items []any) *node {
	m := measurer.Identity()
	for _, item := range items {
		m = measurer.Sum(m, measurer.Measure(item))
	}
	return &node{measurement{measurer, m}, items}
}

func (n *node) String() string {
	var b strings.Builder
	first := true
	b.WriteString("node{")
	for _, i := range n.children {
		if first {
			first = false
		} else {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprint(i))
	}
	b.WriteString("}")
	return b.String()
}

func (n *node) toDigit() *digit {
	return &digit{n._measurement, n.children}
}

func (n *node) Each(f iterFunc) bool {
	for _, item := range n.children {
		if !iterateEach(item, f) {
			return false
		}
	}
	return true
}

func (n *node) EachReverse(f iterFunc) bool {
	for i := len(n.children); i > 0; {
		i--
		if !iterateEachReverse(n.children[i], f) {
			return false
		}
	}
	return true
}

func dup[V any](slice []V) []V {
	result := make([]V, len(slice))
	copy(result, slice)
	return result
}

// Helper function to group an array of elements into an array of nodes.
// m: measurer for nodes
// items: items
// returns array of nodes
func nodes(m measurer, items []any) []*node {
	return nnodes(m, items, []*node{})
}
func nnodes(m measurer, items []any, result []*node) []*node {
	switch len(items) {
	case 2, 3:
		return append(result, newNode(m, items))
	case 4:
		return append(result, newNode(m, dup(items[:2])), newNode(m, dup(items[2:])))
	default:
		result = append(result, newNode(m, dup(items[:3])))
		return nnodes(m, items[3:], result)
	}
}
