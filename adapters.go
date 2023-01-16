package lazyfingertree

import (
	"fmt"
)

type Predicate[M any] func(measure M) bool

type IterFunc[V any] func(value V) bool

type FingerTree[MS Measurer[V, M], V, M any] struct {
	f fingerTree
}

func wrapTree[MS Measurer[V, M], V, M any](tree fingerTree) FingerTree[MS, V, M] {
	return FingerTree[MS, V, M]{tree}
}

var ErrBadValue = fmt.Errorf("%w, bad value", ErrFingerTree)

func wrapPredicate[M any](pred Predicate[M]) func(any) bool {
	return func(m any) bool {
		if wm, ok := m.(M); !ok {
			panic(ErrBadValue)
		} else {
			return pred(wm)
		}
	}
}

func wrapIter[V any](iter IterFunc[V]) func(any) bool {
	return func(v any) bool {
		if wv, ok := v.(V); !ok {
			panic(ErrBadValue)
		} else {
			return iter(wv)
		}
	}
}

func null[T any]() T {
	return *new(T)
}

func (t FingerTree[MS, V, M]) AddFirst(value any) FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](t.f.AddFirst(value))
}

func (t FingerTree[MS, V, M]) AddLast(value any) FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](t.f.AddLast(value))
}

func (t FingerTree[MS, V, M]) RemoveFirst() FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](t.f.RemoveFirst())
}

func (t FingerTree[MS, V, M]) RemoveLast() FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](t.f.RemoveLast())
}
func (t FingerTree[MS, V, M]) PeekFirst() V {
	if cv, ok := t.f.PeekFirst().(V); !ok {
		panic(ErrBadValue)
	} else {
		return cv
	}
}

func (t FingerTree[MS, V, M]) PeekLast() any {
	if cv, ok := t.f.PeekLast().(V); !ok {
		panic(ErrBadValue)
	} else {
		return cv
	}
}

func (t FingerTree[MS, V, M]) Concat(other FingerTree[MS, V, M]) FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](t.f.Concat(other.f))
}

func (t FingerTree[MS, V, M]) Split(predicate Predicate[M]) (FingerTree[MS, V, M], FingerTree[MS, V, M]) {
	left, right := t.f.Split(wrapPredicate(predicate))
	return wrapTree[MS, V, M](left), wrapTree[MS, V, M](right)
}

func (t FingerTree[MS, V, M]) ToSlice() []V {
	s := t.f.ToSlice()
	result := make([]V, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = s[i].(V)
	}
	return result
}

func (t FingerTree[MS, V, M]) diagstr() string {
	return diag(t.f)
}

func (t FingerTree[MS, V, M]) IsEmpty() bool {
	return IsEmpty(t.f)
}

func (t FingerTree[MS, V, M]) Measure() M {
	if cm, ok := Measure(t.f).(M); !ok {
		panic(ErrBadValue)
	} else {
		return cm
	}
}

func (t FingerTree[MS, V, M]) TakeUntil(pred Predicate[M]) FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](TakeUntil(t.f, wrapPredicate(pred)))
}

func (t FingerTree[MS, V, M]) DropUntil(pred Predicate[M]) FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](DropUntil(t.f, wrapPredicate(pred)))
}

func (t FingerTree[MS, V, M]) Each(iter IterFunc[V]) {
	Each(t.f, wrapIter(iter))
}

func (t FingerTree[MS, V, M]) EachReverse(iter IterFunc[V]) {
	EachReverse(t.f, wrapIter(iter))
}

type Measurer[V, M any] interface {
	Identity() M
	// measuring a value could technically produce an error but really should not
	// make sure to validate inputs or to use a panic if you need error support
	Measure(value V) M
	Sum(a M, b M) M
}

type adaptedMeasurer[MS Measurer[V, M], V, M any] struct {
	am MS
}

func (m adaptedMeasurer[MS, V, M]) Identity() any {
	return m.am.Identity()
}

func (m adaptedMeasurer[MS, V, M]) Measure(value any) any {
	if v, ok := value.(V); !ok {
		fmt.Println("Value: ", value)
		panic(fmt.Errorf("%w, wrong value type to measure: %s", ErrBadValue, value))
	} else {
		return m.am.Measure(v)
	}
}

func (m adaptedMeasurer[MS, V, M]) Sum(a any, b any) any {
	if va, ok := a.(M); !ok {
		panic(ErrBadValue)
	} else if vb, ok2 := b.(M); !ok2 {
		panic(ErrBadValue)
	} else {
		return m.am.Sum(va, vb)
	}
}

func FromArray[MS Measurer[V, M], V, M any](measurer MS, values[]V) FingerTree[MS, V, M] {
	cvt := make([]any, len(values))
	for i := 0; i < len(values); i++ {
		cvt[i] = values[i]
	}
	return wrapTree[MS, V, M](fromArray(adaptedMeasurer[MS, V, M]{measurer}, cvt))
}
