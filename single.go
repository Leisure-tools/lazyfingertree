package lazyfingertree

import "fmt"

// A finger-tree which contains exactly one element.
type singleTree struct {
	_measurement measurement
	value        any
}

func newSingleTree(measurer measurer, value any) *singleTree {
	return &singleTree{measurement{measurer, measurer.Measure(value)}, value}
}

func (s *singleTree) measurement() measurement {
	return s._measurement
}

type nodeMeasurer struct {
	identity any
	measurer measurer
}

func (m nodeMeasurer) Identity() any {
	return m.identity
}

func (m nodeMeasurer) Measure(v any) any {
	return asNode(v)._measurement.value
}

func (m nodeMeasurer) Sum(a any, b any) any {
	return m.measurer.Sum(a, b)
}

func makeEmptyMid(m measurer) fingerTree {
	return newEmptyTree(nodeMeasurer{m.Identity(), m})
}

func (s *singleTree) String() string {
	return fmt.Sprintf("singleTree{%v}", s.value)
}

func (s *singleTree) AddFirst(value any) fingerTree {
	m := measurerFor(s)
	return newDeepTree(m,
		newDigit(m, []any{value}),
		makeEmptyMid(m),
		newDigit(m, []any{s.value}),
	)
}

func (s *singleTree) AddLast(value any) fingerTree {
	m := measurerFor(s)
	return newDeepTree(m,
		newDigit(m, []any{s.value}),
		newEmptyTree(m),
		newDigit(m, []any{value}),
	)
}

func (s *singleTree) RemoveFirst() fingerTree {
	return newEmptyTree(measurerFor(s))
}

func (s *singleTree) RemoveLast() fingerTree {
	return newEmptyTree(measurerFor(s))
}

func (s *singleTree) PeekFirst() any {
	return s.value
}

func (s *singleTree) PeekLast() any {
	return s.value
}

func (s *singleTree) Concat(other fingerTree) fingerTree {
	return other.AddFirst(s.value)
}

func (s *singleTree) splitTree(predicate predicate, initial any) (fingerTree, any, fingerTree) {
	return s._measurement.empty(), s.value, s._measurement.empty()
}

func (s *singleTree) Split(predicate predicate) (fingerTree, fingerTree) {
	if predicate(s._measurement.value) {
		return s._measurement.empty(), s
	}
	return s, s._measurement.empty()
}

func (s *singleTree) ToSlice() []any {
	return []any{s.value}
}
