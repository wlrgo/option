package option

import "iter"

// Option represents a value that may or may not be present. The zero value of
// Option is equivalent to [None].
type Option[T any] struct {
	val T
	ok  bool
}

// FromOK returns [Some] of val if err is nil. Otherwise, it returns [None].
func FromOK[T any](val T, err error) Option[T] {
	if err != nil {
		return None[T]()
	}
	return Some(val)
}

// FromPtr returns [Some] of the pointed-to value. If ptr is nil, it returns
// [None].
func FromPtr[T any](ptr *T) Option[T] {
	if ptr == nil {
		return None[T]()
	}

	return Some(*ptr)
}

// None returns an [Option] containing no value.
func None[T any]() Option[T] {
	return Option[T]{ok: false}
}

// Some returns an [Option] containing the value.
func Some[T any](val T) Option[T] {
	return Option[T]{val: val, ok: true}
}

// And returns [None] if o contains no value. Otherwise, it returns other.
//
// other is evaluated before And is called. Use [Option.AndThen] when the
// second [Option] should be computed only if o contains a value.
func (o Option[T]) And[U any](other Option[U]) Option[U] {
	return o.AndThen(func(T) Option[U] { return other })
}

// AndThen returns [None] if o contains no value. Otherwise, it calls f with
// the contained value and returns the result.
//
// The function f is not called when o contains no value.
func (o Option[T]) AndThen[U any](f func(T) Option[U]) Option[U] {
	if o.IsNone() {
		return None[U]()
	}

	return f(o.val)
}

// Expect returns the contained value. It panics with a custom panic message
// provided by msg if o contains no value.
func (o Option[T]) Expect(msg string) T {
	if o.IsNone() {
		panic(msg)
	}

	return o.val
}

// Filter returns o if it contains a value that satisfies f. Otherwise, it
// returns [None].
//
// The function f is not called when o contains no value.
func (o Option[T]) Filter(f func(T) bool) Option[T] {
	return o.AndThen(func(v T) Option[T] {
		if f(v) {
			return Some(v)
		}

		return None[T]()
	})
}

// Get returns the contained value and true. If o contains no value, it returns
// the zero value of T and false.
func (o Option[T]) Get() (T, bool) {
	if o.IsNone() {
		var zero T
		return zero, false
	}

	return o.val, true
}

// GetOrInsert inserts val into o if o contains no value, then returns a
// pointer to the contained value.
//
// The returned pointer aliases o's storage. See [Option.Insert].
//
// val is evaluated before GetOrInsert is called. Use [Option.GetOrInsertWith]
// to compute a fallback only when it is needed.
//
// See also [Option.Insert], which overwrites the contained value even if o
// already contains one.
func (o *Option[T]) GetOrInsert(val T) *T {
	return o.GetOrInsertWith(func() T { return val })
}

// GetOrInsertDefault inserts the zero value of T into o if o contains no
// value, then returns a pointer to the contained value.
//
// The returned pointer aliases o's storage. See [Option.Insert].
func (o *Option[T]) GetOrInsertDefault() *T {
	return o.GetOrInsertWith(func() T { var t T; return t })
}

// GetOrInsertWith inserts a value computed from f into o if o contains no
// value, then returns a pointer to the contained value.
//
// The returned pointer aliases o's storage. See [Option.Insert].
//
// The function f is not called when o contains a value.
func (o *Option[T]) GetOrInsertWith(f func() T) *T {
	if o.IsNone() {
		o.val = f()
		o.ok = true
	}

	return &o.val
}

// Insert inserts val into o, then returns a pointer to the contained value.
//
// If o already contains a value, the old value is overwritten.
//
// The returned pointer aliases o's storage. It is valid until the next
// [Option.Take] or [Option.Replace] on o. Use [Option.Ptr] for a pointer to a
// copy that does not alias o.
//
// Insert must be called on an addressable [Option]. Map elements cannot take
// pointer methods in place.
//
// See also [Option.GetOrInsert], which does not overwrite a contained value.
func (o *Option[T]) Insert(val T) *T {
	o.val = val
	o.ok = true
	return &o.val
}

// Inspect calls f with the contained value if o contains a value.
// It returns o unchanged.
//
// The function f is not called when o contains no value.
func (o Option[T]) Inspect(f func(T)) Option[T] {
	if o.IsSome() {
		f(o.val)
	}

	return o
}

// IsNone reports whether o contains no value.
func (o Option[T]) IsNone() bool {
	return !o.ok
}

// IsNoneOr reports whether o contains no value or the value matches a
// predicate.
//
// The function f is not called when o contains no value.
func (o Option[T]) IsNoneOr(f func(T) bool) bool {
	return o.IsNone() || f(o.val)
}

// IsSome reports whether o contains a value.
func (o Option[T]) IsSome() bool {
	return o.ok
}

// IsSomeAnd reports whether o contains a value and the value matches a
// predicate.
//
// The function f is not called when o contains no value.
func (o Option[T]) IsSomeAnd(f func(T) bool) bool {
	return o.IsSome() && f(o.val)
}

// Map applies f to the contained value and returns the result as an [Option].
// If o contains no value, it returns [None].
//
// The function f is not called when o contains no value.
func (o Option[T]) Map[U any](f func(T) U) Option[U] {
	return o.AndThen(func(v T) Option[U] { return Some(f(v)) })
}

// MapOr applies f to the contained value and returns the result. If o contains
// no value, it returns defaultVal.
//
// defaultVal is evaluated before MapOr is called. Use [Option.MapOrElse] to
// compute a fallback only when it is needed.
//
// The function f is not called when o contains no value.
func (o Option[T]) MapOr[U any](defaultVal U, f func(T) U) U {
	return o.Map(f).UnwrapOr(defaultVal)
}

// MapOrDefault applies f to the contained value and returns the result. If o
// contains no value, it returns the zero value of U.
//
// The function f is not called when o contains no value.
func (o Option[T]) MapOrDefault[U any](f func(T) U) U {
	return o.Map(f).UnwrapOrDefault()
}

// MapOrElse applies someF to the contained value and returns the result. If o
// contains no value, it calls noneF and returns that result.
//
// The function noneF is not called when o contains a value. The function someF
// is not called when o contains no value.
func (o Option[T]) MapOrElse[U any](noneF func() U, someF func(T) U) U {
	return o.Map(someF).UnwrapOrElse(noneF)
}

// OkOr returns the contained value and a nil error. If o contains no value, it
// returns the zero value of T and err.
//
// err should be non-nil when used as a fallback. A nil err on the [None] path
// returns the zero value of T and a nil error, which is indistinguishable from
// [Some] of the zero value.
//
// err is evaluated before OkOr is called. Use [Option.OkOrElse] to compute an
// error only when it is needed.
func (o Option[T]) OkOr(err error) (T, error) {
	return o.OkOrElse(func() error { return err })
}

// OkOrElse returns the contained value and a nil error. If o contains no
// value, it calls errF and returns the zero value of T and that error.
//
// The error returned by errF should be non-nil. See [Option.OkOr].
//
// The function errF is not called when o contains a value.
func (o Option[T]) OkOrElse(errF func() error) (T, error) {
	if o.IsSome() {
		return o.val, nil
	}

	var t T
	return t, errF()
}

// Or returns o if it contains a value. Otherwise, it returns other.
//
// other is evaluated before Or is called. Use [Option.OrElse] when the fallback
// [Option] should be computed only if o contains no value.
func (o Option[T]) Or(other Option[T]) Option[T] {
	return o.OrElse(func() Option[T] { return other })
}

// OrElse returns o if it contains a value. Otherwise, it calls f and returns
// the result.
//
// The function f is not called when o contains a value.
func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}

	return f()
}

// Ptr returns a pointer to a copy of the contained value. If o contains no
// value, it returns nil.
func (o Option[T]) Ptr() *T {
	if o.IsNone() {
		return nil
	}

	t := o.val
	return &t
}

// Reduce combines o and other with f if both contain a value. If only one
// contains a value, that [Option] is returned. If neither contains a value,
// it returns [None].
//
// The function f is not called unless both contain a value.
//
// Unlike [Option.ZipWith], a single contained value is kept instead of
// returning [None].
func (o Option[T]) Reduce(other Option[T], f func(T, T) T) Option[T] {
	return o.ZipWith(other, f).Or(o).Or(other)
}

// Replace replaces the contained value with val and returns the previous
// [Option]. After Replace, o contains val.
func (o *Option[T]) Replace(val T) Option[T] {
	old := *o
	o.Insert(val)
	return old
}

// Seq returns an iterator over the contained value. If o contains no value,
// the iterator is empty.
func (o Option[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.IsSome() && !yield(o.val) {
			return
		}
	}
}

// Take returns the previous [Option] and leaves [None] in its place.
func (o *Option[T]) Take() Option[T] {
	old := *o
	*o = None[T]()
	return old
}

// TakeIf takes the value out of o if f returns true for the contained value.
//
// If o contains a value and f returns true, TakeIf is equivalent to
// [Option.Take]. Otherwise, it leaves o unchanged and returns [None].
//
// f receives a copy of the contained value and cannot mutate o, even when it
// returns false.
//
// The function f is not called when o contains no value.
func (o *Option[T]) TakeIf(f func(T) bool) Option[T] {
	if o.IsNone() || !f(o.val) {
		return None[T]()
	}

	return o.Take()
}

// Unwrap returns the contained value. It panics if o contains no value.
//
// Because this method may panic, its use is generally discouraged. Panics are
// meant for unrecoverable errors, and may abort the entire program. Instead,
// use [Option.UnwrapOr], [Option.UnwrapOrElse], or [Option.UnwrapOrDefault].
func (o Option[T]) Unwrap() T {
	return o.Expect("option: Unwrap called on None")
}

// UnwrapOr returns the contained value or a provided default.
//
// Arguments passed to UnwrapOr are eagerly evaluated; if you are passing the
// result of a function call, it is recommended to use [Option.UnwrapOrElse],
// which is lazily evaluated.
func (o Option[T]) UnwrapOr(defaultVal T) T {
	if o.IsNone() {
		return defaultVal
	}

	return o.val
}

// UnwrapOrDefault returns the contained value. If o contains no value, it
// returns the zero value of T.
func (o Option[T]) UnwrapOrDefault() T {
	if o.IsNone() {
		var defaultVal T
		return defaultVal
	}

	return o.val
}

// UnwrapOrElse returns the contained value. If o contains no value it calls f
// and returns its result. The function f is not called when o contains a
// value.
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.IsNone() {
		return f()
	}

	return o.val
}

// Xor returns o if only o contains a value, or other if only other contains a
// value. It returns [None] otherwise.
func (o Option[T]) Xor(other Option[T]) Option[T] {
	if o.IsSome() == other.IsSome() {
		return None[T]()
	}

	return o.Or(other)
}

// ZipWith applies f to the contained values of o and other and returns the
// result as an [Option]. If either contains no value, it returns [None].
//
// The function f is not called unless both contain a value.
func (o Option[T]) ZipWith[U, V any](
	other Option[U], f func(T, U) V,
) Option[V] {
	return o.AndThen(func(t T) Option[V] {
		return other.Map(func(u U) V { return f(t, u) })
	})
}

// Collect returns [Some] of a slice of every contained value in opts. If any
// element contains no value, it returns [None].
//
// If opts is empty, Collect returns [Some] of an empty slice.
func Collect[T any](opts []Option[T]) Option[[]T] {
	ts := make([]T, len(opts))

	for i, opt := range opts {
		if opt.IsNone() {
			return None[[]T]()
		}
		ts[i] = opt.val
	}

	return Some(ts)
}

// Flatten unwraps an [Option] that contains another [Option].
func Flatten[T any](o Option[Option[T]]) Option[T] {
	return o.AndThen(func(inner Option[T]) Option[T] { return inner })
}
