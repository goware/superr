package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/goware/superr"
)

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
