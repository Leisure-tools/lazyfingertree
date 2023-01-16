package lazyfingertree

type Fake struct {}

func (f *Fake) splitTree(predicate predicate, initial any) (fingerTree, any, fingerTree) {
	return f, nil, f
}

func (f *Fake) measurement() Measurement {
	return *&Measurement{}
}

func (f *Fake) AddFirst(value any) fingerTree {
	return f
}

func (f *Fake) AddLast(value any) fingerTree {
	return f
}

func (f *Fake) RemoveFirst() fingerTree {
	return f
}

func (f *Fake) RemoveLast() fingerTree {
	return f
}

func (f *Fake) PeekFirst() any {
	return nil
}

func (f *Fake) PeekLast() any {
	return nil
}

func (f *Fake) Concat(other fingerTree) fingerTree {
	return f
}

func (f *Fake) Split(predicate predicate) (fingerTree, fingerTree) {
	return f, f
}

func (f *Fake) ToSlice() []any {
	return []any{}
}
