package jsg

import (
	"encoding/json"
	"fmt"
)

// NewObject returns a new json-go node with object type.
func NewObject() Node {
	return &node{
		value: object{},
	}
}

// NewArray returns a new json-go node with array type.
func NewArray() Node {
	return &node{
		value: array{},
	}
}

// New returns the root node for given JSON bytes.
func New(jsonBytes []byte) (Node, error) {
	var value interface{}
	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		return nil, err
	}

	return &node{
		value: value,
	}, nil
}

func (n node) Type() Type {
	switch n.value.(type) {
	case nil:
		return Null
	case string:
		return String
	case float64:
		return Number
	case bool:
		return Boolean
	case object:
		return Object
	case array:
		return Array
	case error:
		return Error
	}

	return Invalid
}

func (n node) Get(p ...interface{}) Node {
	val := n.value

	for _, key := range p {
		switch k := key.(type) {
		case int: // Array item index
			if a, ok := val.(array); ok {
				if k >= 0 && k < len(a) {
					val = a[k]
					continue
				}

				return newError(fmt.Errorf(errorArrayIndexOutOfBounds, k, len(a)))
			}

			return newError(fmt.Errorf(errorIntegerKeyUsedOnNonArray, val, key))
		case string: // Object field key
			if m, ok := val.(object); ok {
				if next, exists := m[k]; exists {
					val = next
					continue
				}

				return newError(fmt.Errorf(errorPathNotFound, key))
			}

			return newError(fmt.Errorf(errorStringKeyUsedOnNonObject, val, key))
		default:
			return newError(fmt.Errorf(errorPathParamType))
		}
	}

	return newValue(val)
}

func (n node) Len() int {
	if a, ok := n.value.(array); ok {
		return len(a)
	}

	return 0
}
