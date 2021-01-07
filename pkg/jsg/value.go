package jsg

import "fmt"

func newValue(v interface{}) Node {
	var nodeType Type

	switch v.(type) {
	case float64:
		nodeType = Number
	case string:
		nodeType = String
	case bool:
		nodeType = Boolean
	}

	return &node{
		nodeType: nodeType,
		value:    v,
	}
}

func (n node) String() string {
	return fmt.Sprint(n.value)
}

func (n node) Number() float64 {
	if num, ok := n.value.(float64); ok {
		return num
	}

	return 0
}

func (n node) Boolean() bool {
	if b, ok := n.value.(bool); ok {
		return b
	}

	return false
}
