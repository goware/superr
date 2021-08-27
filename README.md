superr
======

SupErr -- Go stdlib errors with super powers. Pronounced *super* with a French accent :D

Build a stack of errors compatible with Go stdlib and errors.Is(..).

## Why? & Goals

IMHO, errors in Go are awesome and simple, but often in practice lose meaning and get messy
with many of application-specific messages that are difficult to test for equality.

Typically, as an error is returned up the call stack from `if err != nil { return .., err }`,
it either a.) loses its meaning on the way up b.) creates messy error messages c.) is difficult
to test equality of the error value against a standard set of error values.

`superr` looks to solve those problems by leveraging the new error support from Go.13 of
wrapped errors (via `fmt.Errorf("xxx: %w", err)` and `errors.Is`).

The idea/mindset here is for any application/package you should always start by declaring
a well-defined set of error scenarios as values. This helps you be thoughtful about the
kinds of errors which will occur in your system, but even more importantly it will allow
you to handle these errors throughout your code by having the ability to test the error
emitted by your application against the standard set of errors. 

The stdlib provides most of these capabilities for us, however, the main challenge is you
can't do `fmt.Errorf("%w: %w", ErrBasicAuthError, appErr)`, (this is invalid). This is where
superr comes in, where you would use `superr.New(ErrBasicAuthError, appErr)` which defines
a nested error composed of both `ErrBasicAuthError` and `appErr` values. There are other
benefits superr provides as well -- check out the source and tests for more tricks.


## Example

```go
// Here we supply a well-defined set of Error types used across our application.
var (
  ErrFail         = errors.New("fail") // generic failure
  ErrRequestFail  = errors.New("request fail")
  ErrAuthDeclined = errors.New("authorization declined")
  ErrTimeout      = errors.New("request timed-out")
  ErrDBQuery      = errors.New("db query error")
)

func main() {
  // ..

  _, err := networkFetch("http://..")
  if err != nil {
    if errors.Is(err, ErrRequestFail) {
      // the program execution will land here for any network failure
      log.Fatal(err)
    } else if errors.Is(err, ErrTimeout) {
      // looks like the network timed-out, perhaps we want to retry..?
      // or at least return to user with some time-out code..
    } else {
      fmt.Println("==> generic error:", err)
      os.Exit(1)
    }
  }
}

func networkFetch(url string) (string, error) {
  // some network code.. etc..
  // ..

  err := fmt.Errorf("http failed to fetch %s", url)
  if err != nil {
    // We return the parent error `ErrRequestFail` with the cause of our `err`,
    // where `ErrRequestFail` is from *standard* list of errors, and `err` is
    // our application specific error message.
    //
    // This allows us to easily test for ErrRequestFail errors by the caller
    // but also provides us the underlining application error which caused
    // the error in the first place.
    //
    // superr error types can be nested many levels too (see superr_test.go).
    return "", superr.New(ErrRequestFail, err)
  }

  // will never be reached.. just example
  return "", nil
}
```

## Example 2

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
