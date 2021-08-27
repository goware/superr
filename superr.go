package superr

import (
	"errors"
	"fmt"
)

// New will build a tree of errors, for example, the underlining data structure
// will look like:
//
// Error{err: ErrFail, cause: Error{err: ErrNetwork; cause: Error{err: ErrWee, cause: nil}}}
func New(err error, causes ...error) error {
	stack := &errStack{err: err}
	if len(causes) == 0 {
		return stack
	}
	parent := stack
	for i := 0; i < len(causes); i++ {
		child := &errStack{err: causes[i]}
		parent.cause = child

		if i == len(causes)-1 {
			break
		}
		parent = child
	}
	return stack
}

// Wrap is synonymous with New, in case someone finds that name more familiar
// or intuitive.
func Wrap(err error, causes ...error) error {
	return New(err, causes...)
}

type Error interface {
	Err() error
}

type Cause interface {
	Cause() error
}

type ErrorWithCause interface {
	Error
	Cause
}

var (
	_ Error = &errStack{}
	_ Cause = &errStack{}
)

type errStack struct {
	err   error
	cause error
}

func (e errStack) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s", e.err, e.cause)
	} else {
		return fmt.Sprintf("%s", e.err)
	}
}

func (e errStack) Err() error {
	return e.err
}

func (e errStack) Cause() error {
	return e.cause
}

func (e errStack) Unwrap() error {
	return e.cause
}

func (e errStack) Is(target error) bool {
	if e.err != nil && errors.Is(e.err, target) {
		return true
	}
	if e.cause != nil && errors.Is(e.cause, target) {
		return true
	}
	return false
}

// GetErrorStack returns the nested tree of errors from superr, in the form
// of a flat list.
func GetErrorStack(err error) []error {
	if err == nil {
		return []error{}
	}

	errs := []error{}
	if e, ok := err.(ErrorWithCause); ok {
		errs = append(errs, e.Err())
	} else {
		errs = append(errs, err)
	}

	for {
		unwrap, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}

		werr := unwrap.Unwrap()
		if werr == nil {
			break
		}

		if e, ok := werr.(ErrorWithCause); ok {
			errs = append(errs, e.Err())
		} else {
			errs = append(errs, werr)
		}
		err = werr
	}

	return errs
}
