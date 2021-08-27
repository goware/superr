package superr_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/goware/superr"
)

var (
	ErrFail     = errors.New("fail")
	ErrAppOpps  = errors.New("oops")
	ErrDeclined = errors.New("declined")
)

func TestErrorStack(t *testing.T) {
	err := superr.New(ErrFail, ErrAppOpps, ErrDeclined, fmt.Errorf("nooo.."))
	assert(t, errors.Is(err, ErrFail), "expecting err is ErrFail")
	assert(t, errors.Is(err, ErrAppOpps), "expecting err is ErrAppOpps")
	assert(t, errors.Is(err, ErrDeclined), "expecting err is ErrDeclined")
	assert(t, strings.Contains(err.Error(), "nooo"), "expecting err string to contain 'nooo'")
}

func TestErrorStackDisjointed(t *testing.T) {
	err1 := superr.New(ErrFail, ErrAppOpps)
	err2 := errors.New("something happened")
	err3 := fmt.Errorf("auth fail: %w", ErrDeclined)

	err := superr.New(err1, err2, err3)

	assert(t, errors.Is(err1, ErrFail), "expecting err is ErrFail")
	assert(t, errors.Is(err1, ErrAppOpps), "expecting err is ErrAppOpps")
	assert(t, errors.Is(err, ErrFail), "expecting err is ErrFail")
	assert(t, errors.Is(err, ErrAppOpps), "expecting err is ErrAppOpps")
	assert(t, errors.Is(err, ErrDeclined), "expecting err is ErrDeclined")
	assert(t, errors.Is(err, err2), "expecting err is err2")
	assert(t, strings.Contains(err.Error(), "auth fail"), "expecting err string to contain 'auth fail'")

	fmt.Println("==> example string output:", err)
}

func assert(t *testing.T, cond bool, msg string) {
	if !cond {
		t.Error(msg)
		t.FailNow()
	}
}
