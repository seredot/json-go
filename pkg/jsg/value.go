package jsg

import "fmt"

func newValue(v interface{}) *Node {
	return &Node{
		value: v,
	}
}

// Str returns the string value. If the value is null or node type is not String, returns an empty string.
func (n Node) Str() string {
	v := n.value

	if v == nil {
		return ""
	}

	if _, ok := v.(error); ok {
		return ""
	}

	return fmt.Sprint(n.value)
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
