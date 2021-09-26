package errorwrapping

import (
	"errors"
	"fmt"
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
	fmt.Println(errors.Is(err, &cause{})) // true

	// 以下コードは
	//   if err == io.ErrUnexpectedEOF
	// 今後以下にすべき
	//   if errors.Is(err, io.ErrUnexpectedEOF)
	//
	//   - Comparisons to io.EOF need not be changed, because io.EOF should never be wrapped
}

func TestErrorsAs(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &cause{})
	var causeError *cause
	if ok := errors.As(err, &causeError); ok {
		fmt.Println(causeError) // cause!  err chainからcauseErrorに合う型をさがしてセットする
	}

	// 以下コードは
	//   if e, ok := err.(*os.PathError); ok
	// 今後以下にすべき
	//   var e *os.PathError
	//   if errors.As(err, &e)
}

type customError struct {
	w error
}

func (*customError) Error() string {
	return "customError"
}

func (e *customError) Unwrap() error {
	return e.w
}

func TestCustomError(t *testing.T) {
	err := &customError{w: &cause{}}
	fmt.Println(err) // customError

	fmt.Println(errors.Unwrap(err))       // cause!
	fmt.Println(errors.Is(err, &cause{})) // true
	var causeError *cause
	if ok := errors.As(err, &causeError); ok {
		fmt.Println(causeError) // cause!  err chainからcauseErrorに合う型をさがしてセットする
	}
}
