package jsg

import (
	"encoding/json"
	"fmt"
)

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

func newValue(v interface{}) *Node {
	return &Node{
		value: v,
	}
}

func newError(err error) *Node {
	return &Node{
		value: err,
	}
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

// Get returns a child node in the node tree. p is the list of keys in the path.
// Keys can be either string or int typed.
// string keys can only be used in Object nodes and returns the child node for the key.
// int keys can only be used in Array nodes and returns the child node for the index.
// Example: root.Get("users", 1) returns the second user in users array.
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

// Set sets a value in the node.
// key can be either string or int typed.
// string key can only be used in Object nodes and sets the value of the property of the object.
// int keys can only be used in Array nodes and sets the item of the array with the given index.
// If the length of the array is smaller than the given index, the length of the array is increased
// and the item value is set to that index. In that case newly created empty item indexes are set
// to nil.
// Example: root.Get(users).Set(1, "Jim") returns the second user in users array.
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

// Str returns the string value. If the value is null or node type is not String, returns an empty string.
// TODO: strict string check
func (n Node) Str() string {
	if str, ok := n.value.(string); ok {
		return str
	}

	return ""
}

// Num returns the floating point value. If the value is null or node type is not Number, returns 0.
func (n Node) Num() float64 {
	if num, ok := n.value.(float64); ok {
		return num
	}

	return 0
}

// Bool returns the bool value. If the value is null or node type is not Boolean, returns false.
func (n Node) Bool() bool {
	if b, ok := n.value.(bool); ok {
		return b
	}

	return false
}

// Raw returns the raw value from the JSON parser. Valid types: nil, string, float64, bool, map[string]interface{}, []interface{}.
func (n Node) Raw() interface{} {
	return n.value
}

// Err returns the error if the node type is Error. Otherwise nil.
func (n Node) Err() error {
	if e, ok := n.value.(error); ok {
		return e
	}

	return nil
}

// SerializeIndent(indent string) []btye
// Serialize() []btye
