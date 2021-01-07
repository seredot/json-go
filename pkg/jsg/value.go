package jsg

import "fmt"

func newValue(v interface{}) Node {
	return &node{
		value: v,
	}
}

func (n node) Str() string {
	return fmt.Sprint(n.value)
}

func (n node) Num() float64 {
	if num, ok := n.value.(float64); ok {
		return num
	}

	return 0
}

func (n node) Bool() bool {
	if b, ok := n.value.(bool); ok {
		return b
	}

	return false
}
