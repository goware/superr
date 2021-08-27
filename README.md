superr
======

SupErr -- Go stdlib errors with super powers.

Build a stack of errors compatible with Go stdlib and errors.Is(..)

## Example

```go
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

fmt.Println(err)
// => fail: oops: something happened: auth fail: declined
```

## LICENSE

MIT
