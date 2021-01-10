package jsg

import (
	"encoding/json"
	"fmt"
)

// NewObject returns a new json-go node with object type.
func NewObject() *Node {
	return &Node{
		value: object{},
	}
}

// NewArray returns a new json-go node with array type.
func NewArray() *Node {
	return &Node{
		value: array{},
	}
}

// New returns the root node for given JSON bytes.
func New(jsonBytes []byte) (*Node, error) {
	var value interface{}
	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		return nil, err
	}

	return &Node{
		value: value,
	}, nil
}

func typeOf(v interface{}) Type {
	switch v.(type) {
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

// Type returns the type of the node.
func (n Node) Type() Type {
	return typeOf(n.value)
}

// TODO: doc
func (n Node) Get(p ...interface{}) *Node {
	val := n.value

	for _, key := range p {
		switch k := key.(type) {
		case string: // Object field key
			if m, ok := val.(object); ok {
				if next, exists := m[k]; exists {
					val = next
					continue
				}

				return newError(fmt.Errorf(errorPathNotFound, key))
			}

			return newError(fmt.Errorf(errorStringKeyUsedOnNonObject, val, key))
		case int: // Array item index
			if a, ok := val.(array); ok {
				if k >= 0 && k < len(a) {
					val = a[k]
					continue
				}

				return newError(fmt.Errorf(errorArrayIndexOutOfBounds, k, len(a)))
			}

			return newError(fmt.Errorf(errorIntegerKeyUsedOnNonArray, val, key))
		default:
			return newError(fmt.Errorf(errorPathParamType))
		}
	}

	return newValue(val)
}

// TODO: doc
func (n *Node) Set(key interface{}, val interface{}) error {
	t := typeOf(val)

	if t == Invalid || t == Error {
		return fmt.Errorf(errorInvalidType)
	}

	switch k := key.(type) {
	case string:
		if o, ok := n.value.(object); ok {
			o[k] = val

			return nil
		}

		return fmt.Errorf(errorStringKeyUsedOnNonObject, val, key)
	case int:
		if a, ok := n.value.(array); ok {
			if k < 0 {
				return fmt.Errorf(errorArrayIndexOutOfBounds, k, len(a))
			} else if k >= len(a) {
				n.value = make(array, k+1)
				dest := n.value.(array)
				_ = copy(dest, a)
				dest[k] = val
			} else {
				a[k] = val
			}

			return nil
		}

		return fmt.Errorf(errorIntegerKeyUsedOnNonArray, val, key)
	default:
		return fmt.Errorf(errorPathParamType)
	}
}

// Len returns the length of the array. If the node type is not Array returns 0.
func (n Node) Len() int {
	if a, ok := n.value.(array); ok {
		return len(a)
	}

	return 0
}

// Del deletes from Object or Array.
// If the node type is Object, deletes given string key from the object. If the key is not present Del is no-op.
// If the node type is Array, deletes given integer index from the array. Given index must be inside array bounds.
// Otherwise returns an error.
func (n *Node) Del(key interface{}) error {
	switch k := key.(type) {
	case string:
		if o, ok := n.value.(object); ok {
			delete(o, k)

			return nil
		}

		return fmt.Errorf(errorStringKeyUsedOnNonObject, n.value, key)
	case int:
		if a, ok := n.value.(array); ok {
			if k < 0 {
				return fmt.Errorf(errorArrayIndexOutOfBounds, k, len(a))
			} else if k >= len(a) {
				return fmt.Errorf(errorArrayIndexOutOfBounds, k, len(a))
			}

			copy(a[k:], a[k+1:])
			n.value = a[:len(a)-1]

			return nil
		}

		return fmt.Errorf(errorIntegerKeyUsedOnNonArray, n.value, key)
	default:
		return fmt.Errorf(errorPathParamType)
	}
}

// Raw returns the raw value from the JSON parser. Valid types: nil, string, float64, bool, map[string]interface{}, []interface{}.
func (n Node) Raw() interface{} {
	return n.value
}

// SerializeIndent(indent string) []btye
// Serialize() []btye
