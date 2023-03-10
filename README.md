Lazy Fingertree in Go based on [Qiao's JavaScript version](https://github.com/qiao/fingertree.js/) of [Ralf Hinze's and Ross Paterson's Haskell version](http://www.soi.city.ac.uk/~ross/papers/FingerTree.html)

The public API is parameterized (defined in [adapters.go](adapters.go)):

You provide your own object that supports the Measurer[Value, Measurement] interface. `Values` are in the leaves of the tree and your `Measurer` computes the `Measurements` in the `Measure()` and `Sum()` methods. `Measurements` can be any go objects but they *should be immutable* or there could be trouble. Please see [Ralf Hinze's and Ross Paterson's finger tree paper](http://www.soi.city.ac.uk/~ross/papers/FingerTree.html) (and the [tests](main_test.go)) for more information.

Here's the go doc:

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# lazyfingertree

```go
import "github.com/zot/lazyfingertree"
```

Package lazyfingertree implements parameterized lazy finger trees. See the \[readme\]\(README.md\) for details.

## Index

- [Variables](<#variables>)
- [type FingerTree](<#type-fingertree>)
  - [func Concat[MS Measurer[V, M], V, M any](trees ...FingerTree[MS, V, M]) FingerTree[MS, V, M]](<#func-concat>)
  - [func FromArray[MS Measurer[V, M], V, M any](measurer MS, values []V) FingerTree[MS, V, M]](<#func-fromarray>)
  - [func (t FingerTree[MS, V, M]) AddFirst(value V) FingerTree[MS, V, M]](<#func-fingertreems-v-m-addfirst>)
  - [func (t FingerTree[MS, V, M]) AddLast(value V) FingerTree[MS, V, M]](<#func-fingertreems-v-m-addlast>)
  - [func (t FingerTree[MS, V, M]) Concat(other FingerTree[MS, V, M]) FingerTree[MS, V, M]](<#func-fingertreems-v-m-concat>)
  - [func (t FingerTree[MS, V, M]) DropUntil(pred Predicate[M]) FingerTree[MS, V, M]](<#func-fingertreems-v-m-dropuntil>)
  - [func (t FingerTree[MS, V, M]) Each(iter IterFunc[V])](<#func-fingertreems-v-m-each>)
  - [func (t FingerTree[MS, V, M]) EachReverse(iter IterFunc[V])](<#func-fingertreems-v-m-eachreverse>)
  - [func (t FingerTree[MS, V, M]) IsEmpty() bool](<#func-fingertreems-v-m-isempty>)
  - [func (t FingerTree[MS, V, M]) Measure() M](<#func-fingertreems-v-m-measure>)
  - [func (t FingerTree[MS, V, M]) PeekFirst() V](<#func-fingertreems-v-m-peekfirst>)
  - [func (t FingerTree[MS, V, M]) PeekLast() V](<#func-fingertreems-v-m-peeklast>)
  - [func (t FingerTree[MS, V, M]) RemoveFirst() FingerTree[MS, V, M]](<#func-fingertreems-v-m-removefirst>)
  - [func (t FingerTree[MS, V, M]) RemoveLast() FingerTree[MS, V, M]](<#func-fingertreems-v-m-removelast>)
  - [func (t FingerTree[MS, V, M]) Split(predicate Predicate[M]) (FingerTree[MS, V, M], FingerTree[MS, V, M])](<#func-fingertreems-v-m-split>)
  - [func (t FingerTree[MS, V, M]) String() string](<#func-fingertreems-v-m-string>)
  - [func (t FingerTree[MS, V, M]) TakeUntil(pred Predicate[M]) FingerTree[MS, V, M]](<#func-fingertreems-v-m-takeuntil>)
  - [func (t FingerTree[MS, V, M]) ToSlice() []V](<#func-fingertreems-v-m-toslice>)
- [type IterFunc](<#type-iterfunc>)
- [type Measurer](<#type-measurer>)
  - [func AsMeasurer[V, M any](m any) Measurer[V, M]](<#func-asmeasurer>)
- [type Predicate](<#type-predicate>)


## Variables

```go
var ErrBadMeasurer = fmt.Errorf("%w, bad measurer", ErrFingerTree)
```

```go
var ErrBadValue = fmt.Errorf("%w, bad value", ErrFingerTree)
```

```go
var ErrEmptyTree = fmt.Errorf("%w, empty tree", ErrFingerTree)
```

```go
var ErrExpectedNode = fmt.Errorf("%w, expected a node", ErrFingerTree)
```

```go
var ErrFingerTree = errors.New("finger tree")
```

```go
var ErrUnsupported = fmt.Errorf("%w, unsupported operation", ErrFingerTree)
```

## type [FingerTree](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L20-L22>)

FingerTree is a parameterized wrapper on a low\-level finger tree.

```go
type FingerTree[MS Measurer[Value, Measure], Value, Measure any] struct {
    // contains filtered or unexported fields
}
```

### func [Concat](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L214>)

```go
func Concat[MS Measurer[V, M], V, M any](trees ...FingerTree[MS, V, M]) FingerTree[MS, V, M]
```

### func [FromArray](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L206>)

```go
func FromArray[MS Measurer[V, M], V, M any](measurer MS, values []V) FingerTree[MS, V, M]
```

Create a finger tree. You shouldn't need to provide the type parameters, Go should be able to infer them from your arguments. So you should just be able to say, t := FromArray\(myMeasurer, \[\]Plant\{plant1, plant2\}\)

### func \(FingerTree\[MS, V, M\]\) [AddFirst](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L66>)

```go
func (t FingerTree[MS, V, M]) AddFirst(value V) FingerTree[MS, V, M]
```

Add a value to the start of the tree.

### func \(FingerTree\[MS, V, M\]\) [AddLast](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L71>)

```go
func (t FingerTree[MS, V, M]) AddLast(value V) FingerTree[MS, V, M]
```

Add a value to the and of the tree.

### func \(FingerTree\[MS, V, M\]\) [Concat](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L108>)

```go
func (t FingerTree[MS, V, M]) Concat(other FingerTree[MS, V, M]) FingerTree[MS, V, M]
```

Join two finger trees together

### func \(FingerTree\[MS, V, M\]\) [DropUntil](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L153>)

```go
func (t FingerTree[MS, V, M]) DropUntil(pred Predicate[M]) FingerTree[MS, V, M]
```

Discard all the initial values in the tree that do not satisfy the predicate

### func \(FingerTree\[MS, V, M\]\) [Each](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L158>)

```go
func (t FingerTree[MS, V, M]) Each(iter IterFunc[V])
```

Iterate through the tree starting at the beginning

### func \(FingerTree\[MS, V, M\]\) [EachReverse](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L163>)

```go
func (t FingerTree[MS, V, M]) EachReverse(iter IterFunc[V])
```

Iterate through the tree starting at the end

### func \(FingerTree\[MS, V, M\]\) [IsEmpty](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L134>)

```go
func (t FingerTree[MS, V, M]) IsEmpty() bool
```

Return whether the tree is empty

### func \(FingerTree\[MS, V, M\]\) [Measure](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L139>)

```go
func (t FingerTree[MS, V, M]) Measure() M
```

Return the measure of all the tree's values

### func \(FingerTree\[MS, V, M\]\) [PeekFirst](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L89>)

```go
func (t FingerTree[MS, V, M]) PeekFirst() V
```

Return the first value in the tree. Make sure to test whether the tree is empty because this will panic if it is.

### func \(FingerTree\[MS, V, M\]\) [PeekLast](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L99>)

```go
func (t FingerTree[MS, V, M]) PeekLast() V
```

Return the last value in the tree. Make sure to test whether the tree is empty because this will panic if it is.

### func \(FingerTree\[MS, V, M\]\) [RemoveFirst](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L77>)

```go
func (t FingerTree[MS, V, M]) RemoveFirst() FingerTree[MS, V, M]
```

Remove the first value in the tree. Make sure to test whether the tree is empty because this will panic if it is.

### func \(FingerTree\[MS, V, M\]\) [RemoveLast](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L83>)

```go
func (t FingerTree[MS, V, M]) RemoveLast() FingerTree[MS, V, M]
```

Remove the last value in the tree. Make sure to test whether the tree is empty because this will panic if it is.

### func \(FingerTree\[MS, V, M\]\) [Split](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L114>)

```go
func (t FingerTree[MS, V, M]) Split(predicate Predicate[M]) (FingerTree[MS, V, M], FingerTree[MS, V, M])
```

Split the tree. The first tree is all the starting values that do not satisfy the predicate. The second tree is the first value that satisfies the predicate, followed by the rest of the values.

### func \(FingerTree\[MS, V, M\]\) [String](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L129>)

```go
func (t FingerTree[MS, V, M]) String() string
```

### func \(FingerTree\[MS, V, M\]\) [TakeUntil](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L148>)

```go
func (t FingerTree[MS, V, M]) TakeUntil(pred Predicate[M]) FingerTree[MS, V, M]
```

Return all the initial values in the tree that do not satisfy the predicate

### func \(FingerTree\[MS, V, M\]\) [ToSlice](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L120>)

```go
func (t FingerTree[MS, V, M]) ToSlice() []V
```

Return a slice containing all of the values in the tree

## type [IterFunc](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L17>)

An IterFunc is a function that takes a value and returns true or false. It's used by \[Each\] and \[EachReverse\]. Returning true means to continue iteration. Returning false means to stop.

```go
type IterFunc[V any] func(value V) bool
```

## type [Measurer](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L24-L33>)

```go
type Measurer[Value, Measure any] interface {
    // The "zero" measure
    Identity() Measure
    // Return the measure for a value.
    // Measuring a value could technically produce an error but really should not.
    // Make sure to validate inputs or to use a panic if you need error support.
    Measure(value Value) Measure
    // Add two measures together
    Sum(a Measure, b Measure) Measure
}
```

### func [AsMeasurer](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L168>)

```go
func AsMeasurer[V, M any](m any) Measurer[V, M]
```

The measurer interface

## type [Predicate](<https://github.com/zot/lazyfingertree/blob/main/adapters.go#L12>)

A Predicate is a function that takes a measure and returns true or false. It's used by \[Split\], \[TakeUntil\], and \[DropUntil\].

```go
type Predicate[M any] func(measure M) bool
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
