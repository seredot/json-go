package jsg

// Type is a json-go node type.
type Type int

const (
	// Undefined indicates a json-go undefined node.
	Undefined Type = iota
	// Null indicates a json-go null node.
	Null
	// String indicates a json-go string node.
	String
	// Number indicates a json-go number node.
	Number
	// Boolean indicates a json-go boolean node.
	Boolean
	// Object indicates a json-go object node.
	Object
	// Array indicates a json-go array node.
	Array
	// Error indicates a json-go error node.
	Error
)

type node struct {
	nodeType Type
	value    interface{}
}

type object = map[string]interface{}
type array = []interface{}

// Node implements a json-go node.
type Node interface {
	Type() Type
	Get(path interface{}) Node
	String() string
	Number() float64
	Boolean() bool
}
