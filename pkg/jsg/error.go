package jsg

const (
	errorPathParamType            = "path params must be string or integer"
	errorPathNotFound             = "path element not found: '%v'"
	errorIntegerKeyUsedOnNonArray = "integer key used on non-array of type: %T, index: %v"
	errorStringKeyUsedOnNonObject = "string key used on non-object of type: %T, key: '%v'"
	errorArrayIndexOutOfBounds    = "array index out of bounds: index %d of length: %d"
)

func newError(err error) Node {
	return &node{
		value: err,
	}
}

func (n node) Err() error {
	if e, ok := n.value.(error); ok {
		return e
	}

	return nil
}
