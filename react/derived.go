package react

// Derive creates a new Prop whose value is computed from a source Observable.
// The returned Prop updates whenever the source changes, deduplicating equal values.
// The second return value stops the derivation when called.
func Derive[S any, T comparable](source Observable[S], fn func(S) T) (*Prop[T], func()) {
	p := NewProp(fn(source.Get()))
	unsub := source.OnChange(func(s S) {
		p.Set(fn(s))
	})
	return p, unsub
}
