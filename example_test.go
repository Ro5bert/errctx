package errctx

import (
	"fmt"
)

func doX() error {
	// Try to do X, but it fails...
	return New("could not do X: bad stuff happened")
	// This is essentially equivalent to:
	// return errors.New("could not do X: bad stuff happened")
}

func doY() error {
	//...
	// Doing Y depends on doing X:
	if err := doX(); err != nil {
		return WithCtx(err, "could not do Y")
	}
	//...
	return nil
}

func Example() {
	//...
	if err := doY(); err != nil {
		fmt.Println(err.Error())
	}
	//...

	// Output: could not do Y: could not do X: bad stuff happened
}
