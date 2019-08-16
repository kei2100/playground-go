package errorwrapping

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

type cause struct{}

func (*cause) Error() string {
	return "cause!"
}

func TestFmtErrorfWrap(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &cause{})
	fmt.Println(err)                // wrapped: cause!
	fmt.Println(errors.Unwrap(err)) // cause!
}

func TestErrorsIs(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &cause{})
	log.Println(errors.Is(err, &cause{})) // true
}

func TestErrorsAs(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &cause{})
	var causeError *cause
	if ok := errors.As(err, &causeError); ok {
		log.Println(causeError) // cause!  err chainからcauseErrorに合う型をさがしてセットする
	}
}

// TODO
type CustomError struct {
}
