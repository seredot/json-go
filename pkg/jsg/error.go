package jsg

import "errors"

const (
	errorPathParamType = "path must be string or integer"
)

func newError(serr string) Node {
	return &node{
		value: errors.New(serr),
	}
}

func (n node) Err() error {
	if e, ok := n.value.(error); ok {
		return e
	}

	return nil
}
