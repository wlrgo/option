// Package option provides a generic Option type for representing values that
// may be absent.
//
// It follows the Rust Option API as closely as Go allows, rather than
// reshaping it into an idiomatic Go design. It also provides Go helpers such
// as [FromOK], [FromPtr], [Option.Get], [Option.Ptr], and [Option.OkOr] where
// Rust has no equivalent.
//
// # Construction
//
// The zero value of [Option] contains no value and is equivalent to [None].
// Use [Some], [FromOK], or [FromPtr] to construct an [Option] containing a
// value.
//
// # Methods and functions
//
// Operations on a single [Option] are methods, including those that introduce
// another type, such as [Option.And], [Option.AndThen], [Option.Map], and
// [Option.ZipWith]. Package-level functions are used for construction
// ([Some], [None], [FromOK], [FromPtr]), collecting ([Collect]), flattening a
// nested [Option] ([Flatten]), and operations that need an extra constraint,
// such as [Compare] and [Equal].
//
// # API
//
// Querying: [Option.IsSome], [Option.IsNone], [Option.IsSomeAnd],
// [Option.IsNoneOr].
//
// Extracting: [Option.Expect], [Option.Unwrap], [Option.UnwrapOr],
// [Option.UnwrapOrDefault], [Option.UnwrapOrElse].
//
// Combining: [Option.And], [Option.AndThen], [Option.Or], [Option.OrElse],
// [Option.Xor].
//
// Transforming: [Option.Map], [Option.MapOr], [Option.MapOrDefault],
// [Option.MapOrElse], [Flatten], [Option.ZipWith], [Option.Filter],
// [Option.Inspect], [Option.Reduce].
//
// Comparing: [Compare], [Equal], [Ge], [Gt], [Le], [Lt]. [None] is less than
// any [Some].
//
// Modifying: [Option.Insert], [Option.GetOrInsert],
// [Option.GetOrInsertDefault], [Option.GetOrInsertWith], [Option.Take],
// [Option.TakeIf], [Option.Replace]. These methods require an addressable
// [Option]. Map elements cannot take pointer methods in place.
//
// Iterating: [Option.Seq] yields the contained value if present. [Collect]
// turns a slice of options into an option of a slice.
//
// Converting: [Option.Get], [Option.Ptr], [Option.OkOr], [Option.OkOrElse].
//
// # Evaluation
//
// Combinators that take a fallback value evaluate it before the call:
// [Option.And], [Option.Or], [Option.UnwrapOr], [Option.MapOr], [Option.OkOr],
// and [Option.GetOrInsert]. Use the Else or With variants to compute a fallback
// only when it is needed.
//
// # Mutation
//
// Pointers returned by [Option.Insert] and [Option.GetOrInsert] alias
// internal storage and are invalidated by [Option.Take] and [Option.Replace].
// [Option.Ptr] returns a pointer to a copy. [Option.TakeIf] passes the
// contained value by copy, so the predicate cannot mutate the option.
package option
