package errctx

import (
	"fmt"
)

type recoverableError struct{}

func (recoverableError) Error() string {
	return "recoverableError"
}

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
		// You can inspect the root error and act depending on its type.
		if _, ok := Root(err).(recoverableError); ok {
			// If the error is a recoverableError (it never is), we do not prematurely exit the program.
		} else {
			fmt.Println(err.Error())
			return // Premature program exit.
		}
	}
	//...

	// Output: could not do Y: could not do X: bad stuff happened
}
