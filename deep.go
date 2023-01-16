package lazyfingertree

import "fmt"

// A finger-tree which contains more than one element.
type deepTree struct {
	measured bool
	_measurement Measurement
	left *digit
	mid fingerTree
	right *digit
}

func newDeepTree(measurer measurer, left *digit, mid fingerTree, right *digit) *deepTree {
	return &deepTree{
		false,
		Measurement{measurer, measurer.Identity()},
		left,
		mid,
		right,
	}
}

func (d *deepTree) diagstr() string {
	return fmt.Sprintf("deepTree{%s, %s, %s}", diag(d.left), diag(d.mid), diag(d.right))
}

func (d *deepTree) measurement() Measurement {
	if !d.measured {
		meas := d._measurement.measurer
		d._measurement.value = meas.Sum(
			meas.Sum(d.left._measurement.value, Measure(d.mid)),
			d.right._measurement.value,
		)
		d.measured = true
	}
	return d._measurement
}

func (d *deepTree) AddFirst(v any) fingerTree {
	var meas = measurerFor(d)
	leftItems := d.left.items
	if len(leftItems) == 4 {
		return newDeepTree(
			meas,
			newDigit(meas, []any{v, leftItems[0]}),
			d.mid.AddFirst(newNode(meas, []any{leftItems[1], leftItems[2], leftItems[3]})),
			d.right,
		)
	}
	digits := make([]any, len(leftItems) + 1)
	digits[0] = v
	copy(digits[1:], leftItems)
	return newDeepTree(
		meas,
		newDigit(meas, digits),
		d.mid,
		d.right,
	)
}

func (d *deepTree) AddLast(v any) fingerTree {
	meas := measurerFor(d)
	rightItems := d.right.items
	if d.right.len() == 4 {
		return newDeepTree(
			meas,
			d.left,
			d.mid.AddLast(newNode(meas, []any{rightItems[0], rightItems[1], rightItems[2]})),
			newDigit(meas, []any{rightItems[3], v}),
		)
	}
	digits := make([]any, len(rightItems) + 1)
	copy(digits, rightItems)
	digits[len(rightItems)] = v
	return newDeepTree(
		meas,
		d.left,
		d.mid,
		newDigit(meas, digits),
	)
}

func (d *deepTree) RemoveFirst() fingerTree {
	meas := measurerFor(d)
	if d.left.len() > 1 {
		return newDeepTree(meas, d.left.removeFirst(), d.mid, d.right)
	} else if !IsEmpty(d.mid) {
		newMid := newDelayed(func() fingerTree {return d.mid.RemoveFirst()})
		midFirst := d.mid.PeekFirst()
		return newDeepTree(meas, asNode(midFirst).toDigit(), newMid, d.right)
	} else if d.right.len() == 1 {
		return newSingleTree(meas, d.right.items[0])
	}
	return newDeepTree(meas, d.right.slice(0, 1), d.mid, d.right.slice(1, d.right.len()))
}

func (d *deepTree) RemoveLast() fingerTree {
	meas := measurerFor(d)
	if d.right.len() > 1 {
		return newDeepTree(meas, d.left, d.mid, d.right.removeLast())
	} else if !IsEmpty(d.mid) {
		newMid := newDelayed(func() fingerTree {return d.mid.RemoveLast()})
		last := d.mid.PeekLast()
		return newDeepTree(meas, d.left, newMid, asNode(last).toDigit())
	} else if d.left.len() == 1 {
		return newSingleTree(meas, d.left.items[0])
	}
	return newDeepTree(meas, d.left.slice(0, d.left.len() - 1), d.mid, d.left.slice(d.left.len() - 1, d.left.len()))
}

func (d *deepTree) PeekFirst() any {
	return d.left.peekFirst()
}

func (d *deepTree) PeekLast() any{
	return d.right.peekLast()
}

func (d *deepTree) Concat(other fingerTree) fingerTree {
	other = force(other)
	if IsEmpty(other) {
		return d
	} else if s, ok := other.(*singleTree); ok {
		return d.AddLast(s.value)
	}
	return app3(d, []any{}, other)
}

func (d *deepTree) Split(predicate predicate) (fingerTree, fingerTree) {
	meas := Measure(d)
	measurer := measurerFor(d)
	if predicate(meas) {
		left, mid, right := d.splitTree(predicate, measurer.Identity())
		return left, right.AddFirst(mid)
	}
	return d, newEmptyTree(measurer)
}

func (d *deepTree) ToSlice() []any {
	result := make([]any, 0, 8)
	Each(d, func(value any) bool {
		result = append(result, value)
		return true
	})
	return result
}

// Helper function to split the tree into 3 parts.
// middle value could be 
func (d *deepTree) splitTree(predicate predicate, initial any) (fingerTree, any, fingerTree) {
	meas := measurerFor(d)
	leftMeasure := meas.Sum(initial, d.left._measurement.value)
    // see if the split point is inside the left tree
	if predicate(leftMeasure) {
		left, mid, right := d.left.dsplit(predicate, initial)
		return fromArray(meas, left), mid, deepLeft(meas, right, d.mid, d.right)
	}
	midMeasure := meas.Sum(leftMeasure, Measure(d.mid))
	// see if the split point is inside the mid tree
	if predicate(midMeasure) {
		mleft, mmid, mright := d.mid.splitTree(predicate, leftMeasure)
		left, mid, right := asNode(mmid).toDigit().dsplit(predicate, meas.Sum(leftMeasure, Measure(mleft)))
		return deepRight(meas, d.left, mleft, left),
			mid,
			deepLeft(meas, right, mright, d.right)
	}
	// the split point is in the right tree
	left, mid, right := d.right.dsplit(predicate, midMeasure)
	return deepRight(meas, d.left, d.mid, left),
		mid,
		fromArray(meas, right)
}

func deepLeft(meas measurer, left []any, mid fingerTree, right *digit) fingerTree {
	if len(left) == 0 {
		if IsEmpty(mid) {
			return fromArray(meas, right.items)
		}
		return newDelayed(func() fingerTree {
			first := asNode(mid.PeekFirst()).toDigit()
			rest := mid.RemoveFirst()
			return newDeepTree(meas, first, rest, right)
		})
	}
	return newDeepTree(meas, newDigit(meas, left), mid, right)
}

func deepRight(meas measurer, left *digit, mid fingerTree, right []any) fingerTree {
	if len(right) == 0 {
		if IsEmpty(mid) {
			return fromArray(meas, left.items)
		}
		return newDelayed(func() fingerTree {
			butLast := mid.RemoveLast()
			last := asNode(mid.PeekLast()).toDigit()
			return newDeepTree(meas, left, butLast, last)
		})
	}
	return newDeepTree(meas, left, mid, newDigit(meas, right))
}

// Helper function to concatenate two finger-trees with additional elements
// in between.
// t1: Left finger-tree
// ts: An array of elements in between the two finger-trees
// t2: Right finger-tree
// returns a new FingerTree
func app3(t1 fingerTree, items []any, t2 fingerTree) fingerTree {
	t1 = force(t1)
	t2 = force(t2)
	if IsEmpty(t1) {
		return prependTree(t2, items)
	} else if IsEmpty(t2) {
		return appendTree(t1, items)
	} else if s, ok := t1.(*singleTree); ok {
		return prependTree(t2, items).AddFirst(s.value)
	} else if s, ok := t2.(*singleTree); ok {
		return appendTree(t1, items).AddLast(s.value)
	}
	d1, _ := t1.(*deepTree)
	d2, _ := t2.(*deepTree)
	return newDeepTree(
		measurerFor(d1),
		d1.left,
		newDelayed(func() fingerTree {
			return app3(
				d1.mid,
				nodes(measurerFor(d1), concat(d1.right.items, items, d2.left.items), []any{}),
				d2.mid)
		}),
		d2.right)
}

// Helper function to group an array of elements into an array of nodes.
// m: measurer for nodes
// items: items
// returns array of nodes
func nodes(m measurer, items []any, result []any) []any {
	switch len(items) {
	case 2, 3:
		return append(result, newNode(m, items))
	case 4:
		return append(result, newNode(m, []any{items[0], items[1]}), newNode(m, []any{items[2], items[3]}))
	default:
		return nodes(m, items[3:], append(result, newNode(m, []any{items[0], items[1], items[2]})))
	}
}
