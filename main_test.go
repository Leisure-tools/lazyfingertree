package lazyfingertree

import (
	"testing"
	"runtime/debug"
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

func verifyTree(t *testing.T, tree FingerTree[width[int, int], int, int], start int, length int) {
	if length == 0 {
		if !tree.IsEmpty() {
			t.Log("Expected tree to be empty but it is not")
			debug.PrintStack()
			t.Fail()
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
	t3 := newTree[int](nums...)
	for i := 0; i <= size; i++ {
		left, right := t3.Split(func(w int) bool {
			return w > i
		})
		verifyTree(t, left, 0, i)
		verifyTree(t, right, i, size - i)
	}
}

func TestSimple(t *testing.T) {
	for i := 1; i <= 99; i++ {
		//fmt.Println("Testing tree: ", i)
		testTree(t, i)
	}
}
