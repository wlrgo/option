package option

import "cmp"

// Compare returns -1 if x is less than y, 0 if x equals y, or +1 if x is
// greater than y.
//
// [None] is less than any [Some] value. Two [None] values are equal. If both
// contain values, they are compared with [cmp.Compare].
//
// For floating-point values, [cmp.Compare] treats NaN as equal to NaN. [Equal]
// uses Go ==, so two [Some] NaN values compare equal here and unequal there.
func Compare[T cmp.Ordered](x, y Option[T]) int {
	switch {
	case x.IsNone() && y.IsNone():
		return 0
	case x.IsNone():
		return -1
	case y.IsNone():
		return 1
	default:
		return cmp.Compare(x.val, y.val)
	}
}

// Equal reports whether x and y are equal.
//
// Two [None] values are equal. A [None] is not equal to any [Some]. Two [Some]
// values are equal if their contained values are equal with Go ==.
//
// For floating-point values, two [Some] NaN values are not equal. [Compare]
// uses [cmp.Compare], so those same values compare equal there.
func Equal[T comparable](x, y Option[T]) bool {
	if x.IsNone() || y.IsNone() {
		return x.IsNone() && y.IsNone()
	}

	return x.val == y.val
}

// Ge reports whether x is greater than or equal to y.
//
// See [Compare] for the ordering of [None] and [Some].
func Ge[T cmp.Ordered](x, y Option[T]) bool {
	return Compare(x, y) >= 0
}

// Gt reports whether x is greater than y.
//
// See [Compare] for the ordering of [None] and [Some].
func Gt[T cmp.Ordered](x, y Option[T]) bool {
	return Compare(x, y) > 0
}

// Le reports whether x is less than or equal to y.
//
// See [Compare] for the ordering of [None] and [Some].
func Le[T cmp.Ordered](x, y Option[T]) bool {
	return Compare(x, y) <= 0
}

// Lt reports whether x is less than y.
//
// See [Compare] for the ordering of [None] and [Some].
func Lt[T cmp.Ordered](x, y Option[T]) bool {
	return Compare(x, y) < 0
}
