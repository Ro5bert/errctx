# errctx
errctx is a simple Go library for adding context (in the form of strings) to errors in situations when an error occurs
and is propagated back through a long call chain; each function in the call chain may or may not want to add context as
to where/why/how the error occurred.

errctx produces error strings of the following form:

    <third piece of context>: <second piece of context>: <first piece of context>: <root error message>
    
Yes, it is essentially a glorified way of writing

    err = errors.New(fmt.Sprintf("%s: %s", "context", err.Error()))
    
See the in-source documentation for details. There is also a [contrived example](example_test.go) to take a look at.
