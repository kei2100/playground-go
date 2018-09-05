package test

import (
	"fmt"
	"testing"
)

func ExampleSetup() {
	setupA := func(t *testing.T) func() {
		fmt.Println("setupA")
		return func() {
			fmt.Println("teardownA")
		}
	}

	setupB := func(t *testing.T) func() {
		fmt.Println("setupB")
		return func() {
			fmt.Println("teardownB")
		}
	}

	teardown := Setup(new(testing.T), setupA, setupB)
	teardown()

	// Output:
	// setupA
	// setupB
	// teardownB
	// teardownA
}
