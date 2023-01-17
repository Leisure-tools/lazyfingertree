package lazyfingertree

import (
	"fmt"
	"strings"
)

// A node is a measured container of either 2 or 3 sub-finger-trees.
type node struct {
	_measurement measurement
	items        []any
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
	for _, i := range n.items {
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
	return newDigit(n._measurement.measurer, n.items)
}
