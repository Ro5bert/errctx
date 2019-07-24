package errctx

import (
	"errors"
	"fmt"
	"strings"
)

const joinSequence = ":"

func Root(err error) error {
	if err, ok := err.(Error); ok {
		return err.Root()
	}
	return err
}

func WithCtx(err error, ctx string) Error {
	if err, ok := err.(Error); ok {
		err.AddCtx(ctx)
		return err
	}
	return &ctxerror{err, []string{ctx}}
}

func WithCtxf(err error, format string, a ...interface{}) Error {
	return WithCtx(err, fmt.Sprintf(format, a...))
}

func New(text string) Error {
	return &ctxerror{root: errors.New(text)}
}

func Newf(format string, a ...interface{}) Error {
	return New(fmt.Sprintf(format, a...))
}

type Error interface {
	error
	Root() error
	AddCtx(ctx string)
	AddCtxf(format string, a ...interface{})
}

type ctxerror struct {
	root error
	ctx  []string
}

func (err ctxerror) Root() error {
	return err.root
}

func (err *ctxerror) AddCtx(ctx string) {
	err.ctx = append(err.ctx, ctx)
}

func (err *ctxerror) AddCtxf(format string, a ...interface{}) {
	err.AddCtx(fmt.Sprintf(format, a...))
}

func (err ctxerror) Error() string {
	builder := strings.Builder{}
	last := strings.TrimSpace(err.ctx[len(err.ctx)-1])
	builder.WriteString(last)
	for i := len(err.ctx) - 2; i >= 0; i-- {
		trimmedCtx := strings.TrimSpace(err.ctx[i])
        writeCtxToBuilder(&builder, last, trimmedCtx)
        last = trimmedCtx
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
