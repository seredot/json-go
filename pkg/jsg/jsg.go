package jsg

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

func (n node) Type() Type {
	return n.nodeType
}

func (n node) Get(p string) Node {
	return nil
}
