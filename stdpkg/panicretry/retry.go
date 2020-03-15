package panicretry

import "errors"

// Do func
func Do(fn func() error) error {
	for {
		switch err := wrap(fn); err {
		case errPanic:
			continue
		case nil:
			return nil
		default:
			return err
		}
	}
}

var errPanic = errors.New("panic err")

func wrap(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errPanic
		}
	}()
	return fn()
}
