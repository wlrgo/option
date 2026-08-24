# option

[![Go Reference](https://pkg.go.dev/badge/github.com/wlrgo/option/v2.svg)](https://pkg.go.dev/github.com/wlrgo/option/v2)
[![CI](https://github.com/wlrgo/option/actions/workflows/ci.yml/badge.svg)](https://github.com/wlrgo/option/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/wlrgo/option)](https://github.com/wlrgo/option/blob/main/go.mod)
[![Release](https://img.shields.io/github/v/release/wlrgo/option)](https://github.com/wlrgo/option/releases)
[![License](https://img.shields.io/github/license/wlrgo/option)](LICENSE)

`Option[T]` for values that may be absent.

## Rust API in Go

This is a [wlrgo](https://github.com/wlrgo) package. wlrgo ports Rust
standard-library types to Go and **intentionally copies the Rust API** instead
of reshaping it into an idiomatic Go design. Later packages in the org follow
the same rule.

Names, combinators, and eager vs lazy evaluation follow
[`std::option::Option`](https://doc.rust-lang.org/std/option/enum.Option.html)
as closely as Go generics allow. A few Go-only helpers exist where Rust has no
equivalent: `Get`, `FromOK`, `FromPtr`, `Ptr`, and `OkOr`.

## Install

Requires Go 1.27 or later.

```bash
go get github.com/wlrgo/option/v2
```

## Example

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/wlrgo/option/v2"
)

func main() {
	opt := option.FromOK(strconv.Atoi("2"))
	if v, ok := opt.Get(); ok {
		fmt.Println(v)
	}

	fmt.Println(opt.Map(func(n int) int { return n * 10 }).UnwrapOr(-1))
}
```

The zero value is `None`. `Some` and `None` construct an option; `FromOK` and
`FromPtr` convert from `(T, error)` and `*T`. Combinators such as `Map`, `And`,
and `Or` are methods. `Compare`, `Flatten`, and `Collect` remain package-level
functions.

## API

| Group | Highlights |
| --- | --- |
| Query | `IsSome`, `IsNone`, `IsSomeAnd`, `IsNoneOr` |
| Extract | `Unwrap`, `Expect`, `UnwrapOr`, `Get` |
| Combine | `And`, `AndThen`, `Or`, `OrElse`, `Xor` |
| Transform | `Map`, `Filter`, `Flatten`, `ZipWith`, `Reduce` |
| Compare | `Compare`, `Equal` (`None` < `Some`) |
| Modify | `Insert`, `Take`, `TakeIf`, `Replace` |
| Iterate | `Seq`, `Collect` |
| Convert | `FromOK`, `FromPtr`, `Ptr`, `OkOr` |

See [pkg.go.dev/github.com/wlrgo/option/v2](https://pkg.go.dev/github.com/wlrgo/option/v2)
for the full API and package contract.

## License

[MIT](LICENSE)
