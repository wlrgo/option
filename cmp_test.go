package option_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/option/v2"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name string
		a, b option.Option[int]
		want int
	}{
		{"none none", option.None[int](), option.None[int](), 0},
		{"zero none", option.Option[int]{}, option.None[int](), 0},
		{"none some", option.None[int](), option.Some(1), -1},
		{"some none", option.Some(1), option.None[int](), 1},
		{"some less", option.Some(1), option.Some(2), -1},
		{"some equal", option.Some(2), option.Some(2), 0},
		{"some greater", option.Some(3), option.Some(1), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := option.Compare(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b option.Option[int]
		want bool
	}{
		{"none none", option.None[int](), option.None[int](), true},
		{"zero none", option.Option[int]{}, option.None[int](), true},
		{"none some", option.None[int](), option.Some(0), false},
		{"some none", option.Some(0), option.None[int](), false},
		{"some equal", option.Some(2), option.Some(2), true},
		{"some different", option.Some(2), option.Some(3), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := option.Equal(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEqual_NaN(t *testing.T) {
	nan := option.Some(math.NaN())

	assert.False(t, option.Equal(nan, nan))
	assert.Equal(t, 0, option.Compare(nan, nan))
}

func TestGe(t *testing.T) {
	tests := []struct {
		name string
		a, b option.Option[int]
		want bool
	}{
		{"none none", option.None[int](), option.None[int](), true},
		{"none some", option.None[int](), option.Some(1), false},
		{"some none", option.Some(1), option.None[int](), true},
		{"some less", option.Some(1), option.Some(2), false},
		{"some equal", option.Some(2), option.Some(2), true},
		{"some greater", option.Some(3), option.Some(1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := option.Ge(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGt(t *testing.T) {
	tests := []struct {
		name string
		a, b option.Option[int]
		want bool
	}{
		{"none none", option.None[int](), option.None[int](), false},
		{"none some", option.None[int](), option.Some(1), false},
		{"some none", option.Some(1), option.None[int](), true},
		{"some less", option.Some(1), option.Some(2), false},
		{"some equal", option.Some(2), option.Some(2), false},
		{"some greater", option.Some(3), option.Some(1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := option.Gt(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLe(t *testing.T) {
	tests := []struct {
		name string
		a, b option.Option[int]
		want bool
	}{
		{"none none", option.None[int](), option.None[int](), true},
		{"none some", option.None[int](), option.Some(1), true},
		{"some none", option.Some(1), option.None[int](), false},
		{"some less", option.Some(1), option.Some(2), true},
		{"some equal", option.Some(2), option.Some(2), true},
		{"some greater", option.Some(3), option.Some(1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := option.Le(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLt(t *testing.T) {
	tests := []struct {
		name string
		a, b option.Option[int]
		want bool
	}{
		{"none none", option.None[int](), option.None[int](), false},
		{"none some", option.None[int](), option.Some(1), true},
		{"some none", option.Some(1), option.None[int](), false},
		{"some less", option.Some(1), option.Some(2), true},
		{"some equal", option.Some(2), option.Some(2), false},
		{"some greater", option.Some(3), option.Some(1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := option.Lt(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func ExampleCompare() {
	fmt.Println(option.Compare(option.None[int](), option.Some(1)))
	fmt.Println(option.Compare(option.Some(1), option.None[int]()))
	fmt.Println(option.Compare(option.Some(1), option.Some(2)))
	fmt.Println(option.Compare(option.Some(2), option.Some(1)))
	fmt.Println(option.Compare(option.Some(1), option.Some(1)))
	fmt.Println(option.Compare(option.None[int](), option.None[int]()))

	// Output:
	// -1
	// 1
	// -1
	// 1
	// 0
	// 0
}

func ExampleEqual() {
	fmt.Println(option.Equal(option.Some(2), option.Some(2)))
	fmt.Println(option.Equal(option.Some(2), option.None[int]()))
	fmt.Println(option.Equal(option.None[int](), option.None[int]()))

	// Output:
	// true
	// false
	// true
}

func ExampleGe() {
	fmt.Println(option.Ge(option.Some(2), option.Some(1)))
	fmt.Println(option.Ge(option.Some(1), option.Some(1)))
	fmt.Println(option.Ge(option.None[int](), option.Some(1)))

	// Output:
	// true
	// true
	// false
}

func ExampleGt() {
	fmt.Println(option.Gt(option.Some(2), option.Some(1)))
	fmt.Println(option.Gt(option.Some(1), option.Some(1)))
	fmt.Println(option.Gt(option.Some(1), option.None[int]()))

	// Output:
	// true
	// false
	// true
}

func ExampleLe() {
	fmt.Println(option.Le(option.Some(1), option.Some(2)))
	fmt.Println(option.Le(option.Some(1), option.Some(1)))
	fmt.Println(option.Le(option.Some(1), option.None[int]()))

	// Output:
	// true
	// true
	// false
}

func ExampleLt() {
	fmt.Println(option.Lt(option.Some(1), option.Some(2)))
	fmt.Println(option.Lt(option.Some(1), option.Some(1)))
	fmt.Println(option.Lt(option.None[int](), option.Some(1)))

	// Output:
	// true
	// false
	// true
}
