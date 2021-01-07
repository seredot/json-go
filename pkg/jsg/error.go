package jsg

import "errors"

const (
	errorPathParamType = "path must be string or integer"
)

func newError(serr string) Node {
	return &node{
		nodeType: Error,
		value:    errors.New(serr),
	}
}
