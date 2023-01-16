Lazy Fingertree in Go based on [Qiao's JavaScript version](https://github.com/qiao/fingertree.js/) of [Ralf Hinze's and Ross Paterson's Haskell version](http://www.soi.city.ac.uk/~ross/papers/FingerTree.html)

The public API is generic but the implementation is currently not generic because of some type resolution problems. In the future, if it's actually possible, the parameters could be pushed down into the implementation to make the code cleaner.

[adapters.go](adapters.go) contains the public API:

You provide your own object that supports the Measurer[Value, Measurement] interface. `Values` are in the leaves of the tree and your `Measurer` computes the `Measurements` in the `Measure()` and `Sum()` methods. `Measurements` can be any go objects but they *should be immutable* or there could be trouble. Please see [Ralf Hinze's and Ross Paterson's finger tree paper](http://www.soi.city.ac.uk/~ross/papers/FingerTree.html) (and the test code) for more information.

The measurement interface:

```go
type Measurer[V, M any] interface {
	Identity() M
	// measuring a value could technically produce an error but really should not
	// make sure to validate inputs or to use a panic if you need error support
	Measure(value V) M
	Sum(a M, b M) M
}
```

You create a finger tree with `FromArray`, given a measurer and a slice of values (you shouldn't need to provide the type parameters, Go should be able to infer them from your arguments):

```go
func FromArray[MS Measurer[V, M], V, M any](measurer MS, values[]V) FingerTree[MS, V, M]
```

You can print a tree with the `Diag(tree)` function (if any of your data implements the `diagstr() string` method, it will use that while printing the tree). Example:

```go
t := FromArray(myMeasurer, []Plant{plant1, plant2})
println(Diag(t))
```

`FingerTree[MS Measurerer[V, M], V, M any]` supports these methods:

```go
AddFirst(value any) FingerTree[MS, V, M]
AddLast(value any) FingerTree[MS, V, M]
RemoveFirst() FingerTree[MS, V, M]
RemoveLast() FingerTree[MS, V, M]
PeekFirst() V
PeekLast() V
Concat(other FingerTree[MS, V, M]) FingerTree[MS, V, M]
Split(predicate Predicate[M]) (FingerTree[MS, V, M], FingerTree[MS, V, M])
ToSlice() []V
IsEmpty() bool
Measure() M
TakeUntil(pred Predicate[M]) FingerTree[MS, V, M]
DropUntil(pred Predicate[M]) FingerTree[MS, V, M]
Each(iter IterFunc[V])
EachReverse(iter IterFunc[V])
```

Precicates operate on measurements and IterFuncs operate on values and are defined as:

```go
type Predicate[M any] func(measure M) bool

type IterFunc[V any] func(value V) bool
```
