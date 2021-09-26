package sync

import (
	"errors"
	"strings"
	"testing"
)

func TestMultiErrGroup(t *testing.T) {
	group := MultiErrGroup{}
	group.Go(func() error { return errors.New("error #1") })
	group.Go(func() error { return errors.New("error #2") })

	errs := group.Wait()
	if len(errs) != 2 {
		t.Fatalf("len(errs) got %v, want 2", len(errs))
	}
	msg := errs[0].Error() + ":" + errs[1].Error()
	if !strings.Contains(msg, "error #1") {
		t.Errorf("msg not contains 'error #1': %v", msg)
	}
	if !strings.Contains(msg, "error #2") {
		t.Errorf("msg not contains 'error #2': %v", msg)
	}
}
