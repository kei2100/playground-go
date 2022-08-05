package generics

import "golang.org/x/exp/constraints"

// https://go.dev/blog/intro-generics

// GMin is a generic min(a, b) function. using golang.org/x/exp/constraints package
//
//	Actual definition of constraints.Ordered as follows
//	```
//	type Ordered interface {
//		Integer | Float | ~string
//	}
//	```
func GMin[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// --- ↓ type constraint examples ↓ ---

// [S interface{~[]E}, E interface{}]
// Here S must be a slice type whose element type can be any type.
//
// Because this is a common case, the enclosing interface{} may be omitted for interfaces in constraint position, and we can simply write:
// [S ~[]E, E interface{}]
//
// Go 1.18 introduces a new predeclared identifier any as an alias for the empty interface type. With that, we arrive at this idiomatic code:
// [S ~[]E, E any]
