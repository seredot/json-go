package jsg

import (
	"encoding/json"
)

// NewObject returns a new json-go node with object type.
func NewObject() Node {
	return &node{
		nodeType: Object,
		value:    array{},
	}
}

// NewArray returns a new json-go node with array type.
func NewArray() Node {
	return &node{
		nodeType: Array,
		value:    object{},
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
		nodeType: Object,
		value:    value,
	}, nil
}

func (n node) Type() Type {
	return n.nodeType
}

func (n node) Get(p interface{}) Node {
	switch k := p.(type) {
	case int: // Array item index
		if a, ok := n.value.(array); ok {
			return newValue(a[k])
		}
	case string: // Object field key
		if m, ok := n.value.(object); ok {
			return newValue(m[k])
		}
	default:
		return newError(errorPathParamType)
	}

	return nil
}

func (n node) Len() int {
	if a, ok := n.value.(array); ok {
		return len(a)
	}

	return 0
}
