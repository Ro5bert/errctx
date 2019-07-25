// Package errctx implements a simple method of adding context (in the form of strings) to errors, as one might want to
// do as an error propagates back through a long call chain.
//
// errctx provides the Error interface, which embeds the standard error interface to produce error strings of the
// following form:
//
//     <third piece of context>: <second piece of context>: <first piece of context>: <root error message>
//
// (The angle brackets are not necessarily present; they exist to represent that they and the enclosed text would be
// replaced by the string described by the enclosed text. Obviously, one can have an (almost) arbitrary number of
// context strings.)
//
// Each context string in an Error is added to the Error via WithCtx or WithCtxf. Alternatively, context strings may be
// added directly to Errors via Error.AddCtx or Error.AddCtxf.
//
// The ordinal adjectives in the example error message refer to the order the context strings are added to an Error. So
// the "first" context string added to an Error appears closest to the root error, the "second" context string the
// second closest, etc.
//
// When an Error is implicitly created via WithCtx or WithCtxf, the given error's error string is used as the "root
// error message". When an Error is explicitly created via New or Newf, the given string is used as the "root error
// message".
package errctx

import (
	"errors"
	"fmt"
	"strings"
)

const joinSequence = ":"

// Root returns the "root" error for the given error, err. If err implements Error, then Root returns err.Root();
// otherwise, Root simply returns err.
func Root(err error) error {
	if err, ok := err.(Error); ok {
		return err.Root()
	}
	return err
}

// WithCtx adds the given context string to the given error, err. If err implements Error, then WithCtx adds the context
// string to err (as if via err.AddCtx) and returns err; otherwise, WithCtx constructs a new Error, adds the context
// string to it, then returns it.
func WithCtx(err error, ctx string) Error {
	if err, ok := err.(Error); ok {
		err.AddCtx(ctx)
		return err
	}
	return &ctxerror{err, []string{ctx}}
}

// WithCtxf adds the context string derived from the given format specifier to the given error, err. If err implements
// Error, then WithCtxf adds the context string to err (as if via err.AddCtx) and returns err; otherwise, WithCtx
// constructs a new Error, adds the context string to it, then returns it.
func WithCtxf(err error, format string, a ...interface{}) Error {
	return WithCtx(err, fmt.Sprintf(format, a...))
}

// New constructs a new Error with no context and a root error with the given text.
func New(text string) Error {
	return &ctxerror{root: errors.New(text)}
}

// Newf constructs a new Error with no context and a root error with text derived from the given format specifier.
func Newf(format string, a ...interface{}) Error {
	return New(fmt.Sprintf(format, a...))
}

// Error is an interface that embeds the standard error interface. Implementors of Error return a formatted error string
// from the Error method that agrees with the package documentation.
type Error interface {
	error

	// Root returns the "root" error of this Error; that is, the error on to which all the context strings are
	// prepended.
	Root() error

	// AddCtx adds the given context string to this Error.
	AddCtx(ctx string)

	// AddCtxf adds the context string derived from the given format specifier to this Error.
	AddCtxf(format string, a ...interface{})
}

// ctxerror is the default implementation of Error.
type ctxerror struct {
	root error
	ctx  []string
}

func (err ctxerror) Root() error {
	return err.root
}

func (err *ctxerror) AddCtx(ctx string) {
	err.ctx = append(err.ctx, strings.TrimSpace(ctx))
}

func (err *ctxerror) AddCtxf(format string, a ...interface{}) {
	err.AddCtx(fmt.Sprintf(format, a...))
}

func (err ctxerror) Error() string {
	builder := strings.Builder{}
	last := err.ctx[len(err.ctx)-1]
	builder.WriteString(last)
	for i := len(err.ctx) - 2; i >= 0; i-- {
		writeCtxToBuilder(&builder, last, err.ctx[i])
		last = err.ctx[i]
	}
	writeCtxToBuilder(&builder, last, err.root.Error())
	return builder.String()
}

func writeCtxToBuilder(builder *strings.Builder, last string, ctx string) {
	if !strings.HasSuffix(last, joinSequence) {
		builder.WriteString(joinSequence)
	}
	builder.WriteString(" ")
	builder.WriteString(ctx)
}
