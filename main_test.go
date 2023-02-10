package lazyfingertree

import (
	"runtime/debug"
	"testing"
)

type width[Value any, M int] int

func (w width[Value, M]) Identity() M {
	return 0
}

func (w width[Value, M]) Measure(v Value) M {
	return 1
}

func (w width[Value, M]) Sum(a M, b M) M {
	return a + b
}

func asWidth[Value any](v any) width[Value, int] {
	w, _ := v.(width[Value, int])
	return w
}

func newWidth[Value any]() width[Value, int] {
	return 0
}

func newTree[V any](values ...V) FingerTree[width[V, int], V, int] {
	return FromArray[width[V, int], V, int](newWidth[V](), values)
}

func treeType[V any](t FingerTree[width[V, int], V, int]) string {
	if t.IsEmpty() {
		return "empty"
	} else {
		count := 0
		t.Each(func(v V) bool {
			count++
			if count > 1 {
				return false
			}
			return true
		})
		if count == 1 {
			return "single"
		} else {
			return "deep"
		}
	}
}

func failIfNot(t *testing.T, cond bool) {
	if !cond {
		t.Fail()
	}
}

func failIfErrNow(t *testing.T, err any) {
	if err != nil {
		t.Log(err)
		debug.PrintStack()
		t.FailNow()
	}
}

func same[V any](a []V, b []V) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if any(a[i]) != any(b[i]) {
			return false
		}
	}
	return true
}

func verifyTree(t *testing.T, tree FingerTree[width[int, int], int, int], start int, length int) {
	if length == 0 {
		if !tree.IsEmpty() {
			t.Log("Expected tree to be empty but it is not")
			debug.PrintStack()
			t.FailNow()
		}
		return
	}
	failIfNot(t, tree.PeekFirst() == start)
	for i := 0; i <= length; i++ {
		offset := i + start
		left, right := tree.Split(func(w int) bool {
			return w > i
		})
		failIfNot(t, left.Measure() == i)
		if !right.IsEmpty() {
			failIfNot(t, right.PeekFirst() == offset)
		}
	}
}

func testTree(t *testing.T, size int) {
	nums := make([]int, size)
	for i := 0; i < len(nums); i++ {
		nums[i] = i
	}
	tree := newTree(nums...)
	for i := 0; i <= size; i++ {
		left, right := tree.Split(func(w int) bool {
			return w > i
		})
		verifyTree(t, left, 0, i)
		verifyTree(t, right, i, size-i)
		merged := left.Concat(right)
		failIfNot(t, same(tree.ToSlice(), merged.ToSlice()))
		if i < 2 || i > size-3 {
			continue
		}
		item1 := left.PeekLast()
		left = left.RemoveLast()
		item2 := right.PeekFirst()
		right = right.RemoveFirst()
		tree := right.Concat(newTree(item1, item2)).Concat(left)
		n := make([]int, 0, size)
		n = append(n, nums[i+1:]...)
		n = append(n, nums[i-1], nums[i])
		n = append(n, nums[:i-1]...)
		failIfNot(t, same(tree.ToSlice(), n))
	}
}

func TestMerge(t *testing.T) {
	newTree(1, 2).Concat(newTree(3, 4))
}

func TestSimple(t *testing.T) {
	for i := 1; i <= 99; i++ {
		testTree(t, i)
	}
}
