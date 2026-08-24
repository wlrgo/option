package option_test

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/option/v2"
)

func TestCollect(t *testing.T) {
	tests := []struct {
		name      string
		give      []option.Option[int]
		wantValue []int
		wantNone  bool
	}{
		{"empty", nil, []int{}, false},
		{"all some", []option.Option[int]{option.Some(1), option.Some(2)}, []int{1, 2}, false},
		{"first none", []option.Option[int]{option.None[int](), option.Some(2)}, nil, true},
		{"middle none", []option.Option[int]{option.Some(1), option.None[int](), option.Some(3)}, nil, true},
		{"last none", []option.Option[int]{option.Some(1), option.None[int]()}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := option.Collect(tt.give)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[option.Option[int]]
		wantValue int
		wantNone  bool
	}{
		{"some some", option.Some(option.Some(6)), 6, false},
		{"some none", option.Some(option.None[int]()), 0, true},
		{"none", option.None[option.Option[int]](), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := option.Flatten(tt.give)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestFromOK(t *testing.T) {
	tests := []struct {
		name      string
		v         int
		err       error
		wantValue int
		wantNone  bool
	}{
		{"nil error", 7, nil, 7, false},
		{"error", 7, errors.New("missing"), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := option.FromOK(tt.v, tt.err)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestFromPtr(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got := option.FromPtr[int](nil)
		assert.True(t, got.IsNone())
	})

	t.Run("non-nil", func(t *testing.T) {
		x := 7
		got := option.FromPtr(&x)
		assert.Equal(t, 7, got.Unwrap())

		x = 9
		assert.Equal(t, 7, got.Unwrap())
	})
}

func TestOption_And(t *testing.T) {
	tests := []struct {
		name      string
		a         option.Option[int]
		b         option.Option[string]
		wantValue string
		wantNone  bool
	}{
		{"some none", option.Some(2), option.None[string](), "", true},
		{"none some", option.None[int](), option.Some("foo"), "", true},
		{"some some", option.Some(2), option.Some("foo"), "foo", false},
		{"none none", option.None[int](), option.None[string](), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.And(tt.b)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_AndThen(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[int]
		wantValue string
		wantNone  bool
		wantCalls int
	}{
		{"some", option.Some(2), "4", false, 1},
		{"none", option.None[int](), "", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.AndThen(func(i int) option.Option[string] {
				calls++
				return option.Some(strconv.Itoa(i * i))
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_Expect(t *testing.T) {
	const msg = "expected value"
	tests := []struct {
		name      string
		give      option.Option[int]
		want      int
		wantPanic bool
	}{
		{"none", option.None[int](), 0, true},
		{"some", option.Some(10), 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.PanicsWithValue(t, msg, func() { tt.give.Expect(msg) })
				return
			}
			got := tt.give.Expect(msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOption_Filter(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[int]
		wantValue int
		wantNone  bool
		wantCalls int
	}{
		{"none", option.None[int](), 0, true, 0},
		{"some rejected", option.Some(3), 0, true, 1},
		{"some kept", option.Some(4), 4, false, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.Filter(func(n int) bool {
				calls++
				return n%2 == 0
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_Get(t *testing.T) {
	tests := []struct {
		name     string
		give     option.Option[int]
		want     int
		wantSome bool
	}{
		{"none", option.None[int](), 0, false},
		{"some", option.Some(7), 7, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.give.Get()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantSome, ok)
		})
	}
}

func TestOption_GetOrInsert(t *testing.T) {
	tests := []struct {
		name string
		give option.Option[int]
		val  int
		want int
	}{
		{"none", option.None[int](), 5, 5},
		{"some", option.Some(2), 5, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.give
			got := opt.GetOrInsert(tt.val)

			assert.Equal(t, tt.want, *got)
			assert.Equal(t, tt.want, opt.Unwrap())

			*got = 99
			assert.Equal(t, 99, opt.Unwrap())
		})
	}
}

func TestOption_GetOrInsertDefault(t *testing.T) {
	tests := []struct {
		name string
		give option.Option[int]
		want int
	}{
		{"none", option.None[int](), 0},
		{"some", option.Some(2), 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.give
			got := opt.GetOrInsertDefault()

			assert.Equal(t, tt.want, *got)
			assert.Equal(t, tt.want, opt.Unwrap())

			*got = 99
			assert.Equal(t, 99, opt.Unwrap())
		})
	}
}

func TestOption_GetOrInsertWith(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[int]
		want      int
		wantCalls int
	}{
		{"none", option.None[int](), 5, 1},
		{"some", option.Some(2), 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.give
			calls := 0
			got := opt.GetOrInsertWith(func() int {
				calls++
				return 5
			})

			assert.Equal(t, tt.wantCalls, calls)
			assert.Equal(t, tt.want, *got)
			assert.Equal(t, tt.want, opt.Unwrap())

			*got = 99
			assert.Equal(t, 99, opt.Unwrap())
		})
	}
}

func TestOption_Insert(t *testing.T) {
	tests := []struct {
		name string
		give option.Option[int]
		val  int
		want int
	}{
		{"none", option.None[int](), 1, 1},
		{"some", option.Some(2), 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.give
			got := opt.Insert(tt.val)

			assert.Equal(t, tt.want, *got)
			assert.Equal(t, tt.want, opt.Unwrap())

			*got = 99
			assert.Equal(t, 99, opt.Unwrap())
		})
	}
}

func TestOption_Inspect(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[int]
		wantValue int
		wantNone  bool
		wantCalls int
	}{
		{"some", option.Some(2), 2, false, 1},
		{"none", option.None[int](), 0, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.Inspect(func(v int) {
				calls++
				assert.Equal(t, tt.wantValue, v)
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_IsNone(t *testing.T) {
	tests := []struct {
		name string
		give option.Option[int]
		want bool
	}{
		{"none", option.None[int](), true},
		{"some", option.Some(1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.IsNone()
			assert.Equal(t, tt.want, got, "want: %v, got: %v", tt.want, got)
		})
	}
}

func TestOption_IsNoneOr(t *testing.T) {
	tests := []struct {
		name string
		give option.Option[int]
		pred func(int) bool
		want bool
	}{
		{
			name: "none",
			give: option.None[int](),
			pred: nil,
			want: true,
		},
		{
			name: "some with true predicate",
			give: option.Some(5),
			pred: func(v int) bool { return v > 0 },
			want: true,
		},
		{
			name: "some with false predicate",
			give: option.Some(-1),
			pred: func(v int) bool { return v > 0 },
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.IsNoneOr(tt.pred)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOption_IsSome(t *testing.T) {
	tests := []struct {
		name string
		give option.Option[int]
		want bool
	}{
		{"none", option.None[int](), false},
		{"some", option.Some(1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.IsSome()
			assert.Equal(t, tt.want, got, "want: %v, got: %v", tt.want, got)
		})
	}
}

func TestOption_IsSomeAnd(t *testing.T) {
	tests := []struct {
		name string
		give option.Option[int]
		pred func(int) bool
		want bool
	}{
		{
			name: "none",
			give: option.None[int](),
			pred: nil,
			want: false,
		},
		{
			name: "some with true predicate",
			give: option.Some(5),
			pred: func(v int) bool { return v > 0 },
			want: true,
		},
		{
			name: "some with false predicate",
			give: option.Some(-1),
			pred: func(v int) bool { return v > 0 },
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.IsSomeAnd(tt.pred)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOption_Map(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[string]
		wantValue int
		wantNone  bool
		wantCalls int
	}{
		{"some", option.Some("foo"), 3, false, 1},
		{"none", option.None[string](), 0, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.Map(func(v string) int {
				calls++
				return len(v)
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_MapOr(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[string]
		want      int
		wantCalls int
	}{
		{"some", option.Some("foo"), 3, 1},
		{"none", option.None[string](), 42, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.MapOr(42, func(v string) int {
				calls++
				return len(v)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestOption_MapOrDefault(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[string]
		want      int
		wantCalls int
	}{
		{"some", option.Some("hi"), 2, 1},
		{"none", option.None[string](), 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.MapOrDefault(func(v string) int {
				calls++
				return len(v)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestOption_MapOrElse(t *testing.T) {
	tests := []struct {
		name             string
		give             option.Option[string]
		want             int
		wantMapCalls     int
		wantDefaultCalls int
	}{
		{"some", option.Some("foo"), 3, 1, 0},
		{"none", option.None[string](), 42, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapCalls := 0
			defaultCalls := 0
			got := tt.give.MapOrElse(
				func() int {
					defaultCalls++
					return 42
				},
				func(v string) int {
					mapCalls++
					return len(v)
				},
			)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantMapCalls, mapCalls)
			assert.Equal(t, tt.wantDefaultCalls, defaultCalls)
		})
	}
}

func TestOption_OkOr(t *testing.T) {
	errMissing := errors.New("missing")

	tests := []struct {
		name    string
		give    option.Option[int]
		want    int
		wantErr error
	}{
		{"none", option.None[int](), 0, errMissing},
		{"some", option.Some(7), 7, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.OkOr(errMissing)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}

	t.Run("nil error", func(t *testing.T) {
		got, err := option.None[int]().OkOr(nil)
		assert.Equal(t, 0, got)
		assert.NoError(t, err)
	})
}

func TestOption_OkOrElse(t *testing.T) {
	errMissing := errors.New("missing")

	tests := []struct {
		name      string
		give      option.Option[int]
		want      int
		wantErr   error
		wantCalls int
	}{
		{"none", option.None[int](), 0, errMissing, 1},
		{"some", option.Some(7), 7, nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got, err := tt.give.OkOrElse(func() error {
				calls++
				return errMissing
			})
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestOption_Or(t *testing.T) {
	tests := []struct {
		name      string
		a, b      option.Option[int]
		wantValue int
		wantNone  bool
	}{
		{"some none", option.Some(2), option.None[int](), 2, false},
		{"none some", option.None[int](), option.Some(67), 67, false},
		{"some some", option.Some(2), option.Some(67), 2, false},
		{"none none", option.None[int](), option.None[int](), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_OrElse(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[int]
		wantValue int
		wantNone  bool
		wantCalls int
	}{
		{"some", option.Some(2), 2, false, 0},
		{"none", option.None[int](), 67, false, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.OrElse(func() option.Option[int] {
				calls++
				return option.Some(67)
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_Ptr(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		assert.Nil(t, option.None[int]().Ptr())
	})

	t.Run("some", func(t *testing.T) {
		opt := option.Some(7)
		p := opt.Ptr()
		assert.Equal(t, 7, *p)

		*p = 9
		assert.Equal(t, 7, opt.Unwrap())
	})
}

func TestOption_Reduce(t *testing.T) {
	tests := []struct {
		name      string
		opt       option.Option[int]
		other     option.Option[int]
		wantValue int
		wantNone  bool
		wantCalls int
	}{
		{"some some", option.Some(12), option.Some(17), 29, false, 1},
		{"some none", option.Some(12), option.None[int](), 12, false, 0},
		{"none some", option.None[int](), option.Some(17), 17, false, 0},
		{"none none", option.None[int](), option.None[int](), 0, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.opt.Reduce(tt.other, func(a, b int) int {
				calls++
				return a + b
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_Replace(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[int]
		val       int
		wantValue int
		wantNone  bool
	}{
		{"none", option.None[int](), 3, 0, true},
		{"some", option.Some(2), 5, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.give
			got := opt.Replace(tt.val)

			assert.Equal(t, tt.val, opt.Unwrap())

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_Seq(t *testing.T) {
	tests := []struct {
		name string
		give option.Option[int]
		want []int
	}{
		{"none", option.None[int](), nil},
		{"some", option.Some(2), []int{2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			for v := range tt.give.Seq() {
				got = append(got, v)
			}
			assert.Equal(t, tt.want, got)
		})
	}

	t.Run("some break", func(t *testing.T) {
		n := 0
		for range option.Some(2).Seq() {
			n++
			break
		}
		assert.Equal(t, 1, n)
	})
}

func TestOption_Take(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[int]
		wantValue int
		wantNone  bool
	}{
		{"none", option.None[int](), 0, true},
		{"some", option.Some(2), 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.give
			got := opt.Take()

			assert.True(t, opt.IsNone())

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_Take_leftoverPointer(t *testing.T) {
	opt := option.None[int]()
	p := opt.Insert(7)
	opt.Take()
	*p = 99

	got, ok := opt.Get()
	assert.Equal(t, 0, got)
	assert.False(t, ok)
	assert.True(t, option.Equal(opt, option.None[int]()))
}

func TestOption_TakeIf(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[int]
		wantValue int
		wantNone  bool
		leftValue int
		leftNone  bool
		wantCalls int
	}{
		{"none", option.None[int](), 0, true, 0, true, 0},
		{"some taken", option.Some(4), 4, false, 0, true, 1},
		{"some kept", option.Some(3), 0, true, 3, false, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.give
			calls := 0
			got := opt.TakeIf(func(n int) bool {
				calls++
				return n%2 == 0
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantNone {
				assert.True(t, got.IsNone())
			} else {
				assert.Equal(t, tt.wantValue, got.Unwrap())
			}

			if tt.leftNone {
				assert.True(t, opt.IsNone())
				return
			}

			assert.Equal(t, tt.leftValue, opt.Unwrap())
		})
	}
}

func TestOption_Unwrap(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[int]
		want      int
		wantPanic bool
	}{
		{"none", option.None[int](), 0, true},
		{"some", option.Some(10), 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() { tt.give.Unwrap() })
				return
			}
			got := tt.give.Unwrap()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOption_UnwrapOr(t *testing.T) {
	tests := []struct {
		name string
		give option.Option[int]
		want int
	}{
		{"none", option.None[int](), 67},
		{"some", option.Some(10), 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.UnwrapOr(67)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOption_UnwrapOrDefault(t *testing.T) {
	tests := []struct {
		name string
		give option.Option[int]
		want int
	}{
		{"none", option.None[int](), 0},
		{"some", option.Some(10), 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.UnwrapOrDefault()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOption_UnwrapOrElse(t *testing.T) {
	tests := []struct {
		name      string
		give      option.Option[int]
		want      int
		wantCalls int
	}{
		{"none", option.None[int](), 67, 1},
		{"some", option.Some(10), 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.UnwrapOrElse(func() int {
				calls++
				return 67
			})
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestOption_Xor(t *testing.T) {
	tests := []struct {
		name      string
		a, b      option.Option[int]
		wantValue int
		wantNone  bool
	}{
		{"some none", option.Some(2), option.None[int](), 2, false},
		{"none some", option.None[int](), option.Some(67), 67, false},
		{"some some", option.Some(2), option.Some(67), 0, true},
		{"none none", option.None[int](), option.None[int](), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Xor(tt.b)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestOption_ZipWith(t *testing.T) {
	tests := []struct {
		name      string
		opt       option.Option[int]
		other     option.Option[string]
		wantValue int
		wantNone  bool
		wantCalls int
	}{
		{"some some", option.Some(1), option.Some("hi"), 3, false, 1},
		{"some none", option.Some(1), option.None[string](), 0, true, 0},
		{"none some", option.None[int](), option.Some("hi"), 0, true, 0},
		{"none none", option.None[int](), option.None[string](), 0, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.opt.ZipWith(tt.other, func(n int, s string) int {
				calls++
				return n + len(s)
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantNone {
				assert.True(t, got.IsNone())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func ExampleCollect() {
	fmt.Println(
		option.Collect([]option.Option[int]{option.Some(1), option.Some(2)}).
			UnwrapOr([]int{-1}),
	)
	fmt.Println(
		option.Collect([]option.Option[int]{option.Some(1), option.None[int]()}).
			UnwrapOr([]int{-1}),
	)

	// Output:
	// [1 2]
	// [-1]
}

func ExampleFlatten() {
	x := option.Some(option.Some(6))
	fmt.Println(option.Flatten(x).UnwrapOr(-1))

	x = option.Some(option.None[int]())
	fmt.Println(option.Flatten(x).UnwrapOr(-1))

	x = option.None[option.Option[int]]()
	fmt.Println(option.Flatten(x).UnwrapOr(-1))

	// Output:
	// 6
	// -1
	// -1
}

func ExampleFromOK() {
	parse := func(s string) (int, error) {
		if s == "7" {
			return 7, nil
		}
		return 0, errors.New("invalid")
	}

	fmt.Println(option.FromOK(parse("7")).UnwrapOr(-1))
	fmt.Println(option.FromOK(parse("x")).UnwrapOr(-1))

	// Output:
	// 7
	// -1
}

func ExampleFromPtr() {
	x := 7
	fmt.Println(option.FromPtr(&x).UnwrapOr(-1))
	fmt.Println(option.FromPtr[int](nil).UnwrapOr(-1))

	// Output:
	// 7
	// -1
}

func ExampleNone() {
	opt := option.None[int]()

	fmt.Println(opt.IsNone())
	// Output: true
}

func ExampleOption() {
	var opt option.Option[int]

	fmt.Println(opt.IsNone())
	// Output: true
}

func ExampleOption_And() {
	x := option.Some(2)
	y := option.None[string]()
	fmt.Println(x.And(y).UnwrapOr("<none>"))

	x = option.None[int]()
	y = option.Some("foo")
	fmt.Println(x.And(y).UnwrapOr("<none>"))

	x = option.Some(2)
	y = option.Some("foo")
	fmt.Println(x.And(y).UnwrapOr("<none>"))

	x = option.None[int]()
	y = option.None[string]()
	fmt.Println(x.And(y).UnwrapOr("<none>"))

	// Output:
	// <none>
	// <none>
	// foo
	// <none>
}

func ExampleOption_AndThen() {
	sqThenToString := func(x uint32) option.Option[string] {
		if x != 0 && x > math.MaxUint32/x {
			return option.None[string]()
		}

		return option.Some(strconv.FormatUint(uint64(x*x), 10))
	}

	fmt.Println(option.Some[uint32](2).AndThen(sqThenToString).UnwrapOr("<none>"))
	fmt.Println(
		option.Some[uint32](1_000_000).AndThen(sqThenToString).UnwrapOr("<none>"),
	)
	fmt.Println(option.None[uint32]().AndThen(sqThenToString).UnwrapOr("<none>"))

	// Output:
	// 4
	// <none>
	// <none>
}

func ExampleOption_Expect() {
	defer func() { fmt.Println(recover()) }()

	fmt.Println(option.Some("value").Expect("fruits are healthy"))
	option.None[string]().Expect("string should not be empty")

	// Output:
	// value
	// string should not be empty
}

func ExampleOption_Filter() {
	isEven := func(n int) bool {
		return n%2 == 0
	}

	fmt.Println(option.None[int]().Filter(isEven).UnwrapOr(-1))
	fmt.Println(option.Some(3).Filter(isEven).UnwrapOr(-1))
	fmt.Println(option.Some(4).Filter(isEven).UnwrapOr(-1))

	// Output:
	// -1
	// -1
	// 4
}

func ExampleOption_Get() {
	if v, ok := option.Some(7).Get(); ok {
		fmt.Println(v)
	}

	if _, ok := option.None[int]().Get(); !ok {
		fmt.Println("none")
	}

	// Output:
	// 7
	// none
}

func ExampleOption_GetOrInsert() {
	opt := option.None[int]()
	p := opt.GetOrInsert(5)
	fmt.Println(*p)

	*p = 7
	fmt.Println(opt.Unwrap())
	fmt.Println(*opt.GetOrInsert(9))

	// Output:
	// 5
	// 7
	// 7
}

func ExampleOption_GetOrInsertDefault() {
	opt := option.None[int]()
	p := opt.GetOrInsertDefault()
	fmt.Println(*p)

	*p = 7
	fmt.Println(opt.Unwrap())

	// Output:
	// 0
	// 7
}

func ExampleOption_GetOrInsertWith() {
	opt := option.None[int]()
	p := opt.GetOrInsertWith(func() int { return 5 })
	fmt.Println(*p)

	*p = 7
	fmt.Println(opt.Unwrap())

	// Output:
	// 5
	// 7
}

func ExampleOption_Insert() {
	opt := option.None[int]()
	p := opt.Insert(1)
	fmt.Println(*p)
	fmt.Println(opt.Unwrap())

	*p = 2
	fmt.Println(opt.Unwrap())

	// Output:
	// 1
	// 1
	// 2
}

func ExampleOption_Inspect() {
	x := option.Some(2).Inspect(func(v int) { fmt.Println("got:", v) })

	fmt.Println(x.Unwrap())

	option.None[int]().Inspect(func(v int) { fmt.Println("got:", v) })

	// Output:
	// got: 2
	// 2
}

func ExampleOption_IsNone() {
	fmt.Println(option.None[int]().IsNone())
	fmt.Println(option.Some(2).IsNone())

	// Output:
	// true
	// false
}

func ExampleOption_IsNoneOr() {
	positive := func(v int) bool { return v > 0 }

	fmt.Println(option.None[int]().IsNoneOr(positive))
	fmt.Println(option.Some(0).IsNoneOr(positive))

	// Output:
	// true
	// false
}

func ExampleOption_IsSome() {
	fmt.Println(option.None[int]().IsSome())
	fmt.Println(option.Some(2).IsSome())

	// Output:
	// false
	// true
}

func ExampleOption_IsSomeAnd() {
	positive := func(v int) bool { return v > 0 }

	fmt.Println(option.None[int]().IsSomeAnd(positive))
	fmt.Println(option.Some(5).IsSomeAnd(positive))

	// Output:
	// false
	// true
}

func ExampleOption_Map() {
	x := option.Some("Hello, World!")
	fmt.Println(x.Map(func(v string) int { return len(v) }).UnwrapOr(-1))

	y := option.None[string]()
	fmt.Println(y.Map(func(v string) int { return len(v) }).UnwrapOr(-1))

	// Output:
	// 13
	// -1
}

func ExampleOption_MapOr() {
	x := option.Some("foo")
	fmt.Println(x.MapOr(42, func(v string) int { return len(v) }))

	x = option.None[string]()
	fmt.Println(x.MapOr(42, func(v string) int { return len(v) }))

	// Output:
	// 3
	// 42
}

func ExampleOption_MapOrDefault() {
	x := option.Some("hi")
	y := option.None[string]()

	fmt.Println(x.MapOrDefault(func(v string) int { return len(v) }))
	fmt.Println(y.MapOrDefault(func(v string) int { return len(v) }))

	// Output:
	// 2
	// 0
}

func ExampleOption_MapOrElse() {
	i := 21

	x := option.Some("foo")
	fmt.Println(
		x.MapOrElse(
			func() int { return 2 * i },
			func(v string) int { return len(v) },
		),
	)

	x = option.None[string]()
	fmt.Println(
		x.MapOrElse(
			func() int { return 2 * i },
			func(v string) int { return len(v) },
		),
	)

	// Output:
	// 3
	// 42
}

func ExampleOption_OkOr() {
	fmt.Println(option.Some(7).OkOr(errors.New("missing")))
	fmt.Println(option.None[int]().OkOr(errors.New("missing")))

	// Output:
	// 7 <nil>
	// 0 missing
}

func ExampleOption_OkOrElse() {
	fmt.Println(option.Some(7).OkOrElse(func() error { return errors.New("missing") }))
	fmt.Println(option.None[int]().OkOrElse(func() error { return errors.New("missing") }))

	// Output:
	// 7 <nil>
	// 0 missing
}

func ExampleOption_Or() {
	x := option.Some(2)
	y := option.None[int]()
	fmt.Println(x.Or(y).UnwrapOr(-1))

	x = option.None[int]()
	y = option.Some(100)
	fmt.Println(x.Or(y).UnwrapOr(-1))

	x = option.Some(2)
	y = option.Some(100)
	fmt.Println(x.Or(y).UnwrapOr(-1))

	x = option.None[int]()
	y = option.None[int]()
	fmt.Println(x.Or(y).UnwrapOr(-1))

	// Output:
	// 2
	// 100
	// 2
	// -1
}

func ExampleOption_OrElse() {
	nobody := func() option.Option[string] { return option.None[string]() }
	vikings := func() option.Option[string] { return option.Some("vikings") }

	fmt.Println(option.Some("barbarians").OrElse(vikings).UnwrapOr("<none>"))
	fmt.Println(option.None[string]().OrElse(vikings).UnwrapOr("<none>"))
	fmt.Println(option.None[string]().OrElse(nobody).UnwrapOr("<none>"))

	// Output:
	// barbarians
	// vikings
	// <none>
}

func ExampleOption_Ptr() {
	p := option.Some(7).Ptr()
	fmt.Println(*p)

	fmt.Println(option.None[int]().Ptr() == nil)

	// Output:
	// 7
	// true
}

func ExampleOption_Reduce() {
	add := func(a, b int) int { return a + b }

	fmt.Println(option.Some(12).Reduce(option.Some(17), add).UnwrapOr(-1))
	fmt.Println(option.Some(12).Reduce(option.None[int](), add).UnwrapOr(-1))
	fmt.Println(option.None[int]().Reduce(option.Some(17), add).UnwrapOr(-1))
	fmt.Println(option.None[int]().Reduce(option.None[int](), add).UnwrapOr(-1))

	// Output:
	// 29
	// 12
	// 17
	// -1
}

func ExampleOption_Replace() {
	opt := option.Some(2)
	fmt.Println(opt.Replace(5).UnwrapOr(-1))
	fmt.Println(opt.Unwrap())

	opt = option.None[int]()
	fmt.Println(opt.Replace(3).UnwrapOr(-1))
	fmt.Println(opt.Unwrap())

	// Output:
	// 2
	// 5
	// -1
	// 3
}

func ExampleOption_Seq() {
	for v := range option.Some(2).Seq() {
		fmt.Println(v)
	}

	for range option.None[int]().Seq() {
		fmt.Println("none")
	}

	// Output:
	// 2
}

func ExampleOption_Take() {
	opt := option.Some(2)
	fmt.Println(opt.Take().UnwrapOr(-1))
	fmt.Println(opt.UnwrapOr(-1))

	opt = option.None[int]()
	fmt.Println(opt.Take().UnwrapOr(-1))
	fmt.Println(opt.UnwrapOr(-1))

	// Output:
	// 2
	// -1
	// -1
	// -1
}

func ExampleOption_TakeIf() {
	isEven := func(n int) bool { return n%2 == 0 }

	opt := option.Some(4)
	fmt.Println(opt.TakeIf(isEven).UnwrapOr(-1))
	fmt.Println(opt.UnwrapOr(-1))

	opt = option.Some(3)
	fmt.Println(opt.TakeIf(isEven).UnwrapOr(-1))
	fmt.Println(opt.UnwrapOr(-1))

	opt = option.None[int]()
	fmt.Println(opt.TakeIf(isEven).UnwrapOr(-1))

	// Output:
	// 4
	// -1
	// -1
	// 3
	// -1
}

func ExampleOption_Unwrap() {
	defer func() { fmt.Println(recover()) }()

	fmt.Println(option.Some("air").Unwrap())
	option.None[string]().Unwrap()

	// Output:
	// air
	// option: Unwrap called on None
}

func ExampleOption_UnwrapOr() {
	fmt.Println(option.None[string]().UnwrapOr("bike"))
	fmt.Println(option.Some("car").UnwrapOr("bike"))

	// Output:
	// bike
	// car
}

func ExampleOption_UnwrapOrDefault() {
	fmt.Println(option.None[int]().UnwrapOrDefault())
	fmt.Println(option.Some(123).UnwrapOrDefault())

	// Output:
	// 0
	// 123
}

func ExampleOption_UnwrapOrElse() {
	x := 5
	fmt.Println(option.None[int]().UnwrapOrElse(func() int { return 2 * x }))
	fmt.Println(option.Some(6).UnwrapOrElse(func() int { return 2 * x }))

	// Output:
	// 10
	// 6
}

func ExampleOption_Xor() {
	x := option.Some(2)
	y := option.None[int]()
	fmt.Println(x.Xor(y).UnwrapOr(-1))

	x = option.None[int]()
	y = option.Some(2)
	fmt.Println(x.Xor(y).UnwrapOr(-1))

	x = option.Some(2)
	y = option.Some(22)
	fmt.Println(x.Xor(y).UnwrapOr(-1))

	x = option.None[int]()
	y = option.None[int]()
	fmt.Println(x.Xor(y).UnwrapOr(-1))

	// Output:
	// 2
	// 2
	// -1
	// -1
}

func ExampleOption_ZipWith() {
	x := option.Some(1)
	y := option.Some("hi")
	fmt.Println(x.ZipWith(y, func(n int, s string) int { return n + len(s) }).UnwrapOr(-1))
	fmt.Println(x.ZipWith(option.None[string](), func(n int, s string) int { return n + len(s) }).UnwrapOr(-1))

	// Output:
	// 3
	// -1
}

func ExampleSome() {
	opt := option.Some(10)

	fmt.Println(opt.IsSome())
	fmt.Println(opt.Unwrap())
	// Output:
	// true
	// 10
}
