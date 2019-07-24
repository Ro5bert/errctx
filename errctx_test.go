package errctx

import (
	"errors"
	"testing"
)

func TestRoot(t *testing.T) {
	err := errors.New("bad stuff")
	if Root(err) != err {
		t.Fail()
	}
	if Root(WithCtx(err, "ctx")) != err {
		t.Fail()
	}
}

func TestCtxFormatting(t *testing.T) {
	err := errors.New("oh no")
	err = WithCtx(err, "1")
	err = WithCtx(err, "       2:                   ")
	err = WithCtx(err, "3:")
	err = WithCtx(err, "  4        ")
	if err.Error() != "4: 3: 2: 1: oh no" {
		t.Errorf("output: %s", err.Error())
	}
}
