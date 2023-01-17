package lazyfingertree

import (
	"fmt"
	"strings"
)

// A digit is a measured container of one to four elements.
// this is not a FingerTree, it only shares some of the methods
type digit struct {
	_measurement measurement
	items        []any
}

func newDigit(measurer measurer, items []any) *digit {
	m := measurer.Identity()
	for _, item := range items {
		m = measurer.Sum(m, measurer.Measure(item))
	}
	return &digit{measurement{measurer, m}, items}
}

func (d *digit) String() string {
	var b strings.Builder
	first := true
	b.WriteString("digit{")
	for _, i := range d.items {
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

func (d *digit) len() int {
	return len(d.items)
}

func (d *digit) getMeasurement() measurement {
	return d._measurement
}

func (d *digit) removeFirst() *digit {
	return d.slice(1, len(d.items))
}

func (d *digit) removeLast() *digit {
	return d.slice(0, len(d.items)-1)
}

func (d *digit) slice(start int, end int) *digit {
	// don't keep large, partly-used slices around
	if end-start > 10 && cap(d.items)/2 > end-start {
		cpy := make([]any, start-end)
		copy(cpy, d.items[start:end])
		return newDigit(d._measurement.measurer, cpy)
	}
	return newDigit(d._measurement.measurer, d.items[start:end])
}

func (d *digit) peekFirst() any {
	return d.items[0]
}

func (d *digit) peekLast() any {
	return d.items[len(d.items)-1]
}

// Split the digit into 3 parts, in which the left part is the elements
// that does not satisfy the predicate, the middle part is the first
// element that satisfies the predicate and the last part is the rest
// elements.
func (d *digit) dsplit(predicate predicate, initial any) ([]any, any, []any) {
	if len(d.items) == 1 {
		return []any{}, d.items[0], []any{}
	}
	m := initial
	i := 0
	var item any = nil
	meas := d._measurement.measurer
	for i, item = range d.items {
		m = meas.Sum(m, meas.Measure(item))
		if predicate(m) {
			break
		}
	}
	return d.items[0:i], item, d.items[i+1:]
}
