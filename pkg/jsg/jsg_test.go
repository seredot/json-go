package jsg

import (
	"testing"

	"gotest.tools/assert"
)

func TestJSONObject(t *testing.T) {
	object, err := New([]byte(`{
		"string": "foo",
		"number": 123.4,
		"boolean": true
	}`))

	// Parsing
	assert.NilError(t, err)

	// String value
	assert.Equal(t, "foo", object.Get("string").Str())

	// Number value
	assert.Equal(t, 123.4, object.Get("number").Num())

	// Boolean value
	assert.Equal(t, true, object.Get("boolean").Bool())
}

func TestJSONArray(t *testing.T) {
	array, err := New([]byte(`[
		1, 2, 3	
	]`))

	// Parsing
	assert.NilError(t, err)

	// Item values
	assert.Equal(t, 1.0, array.Get(0).Num())
	assert.Equal(t, 2.0, array.Get(1).Num())
	assert.Equal(t, 3.0, array.Get(2).Num())

	// Array length
	assert.Equal(t, 3, array.Len())

	array, err = New([]byte(`[
		"a", "b"
	]`))

	// Parsing
	assert.NilError(t, err)

	// Item values
	assert.Equal(t, "a", array.Get(0).Str())
	assert.Equal(t, "b", array.Get(1).Str())

	// Array length
	assert.Equal(t, 2, array.Len())

	array, err = New([]byte(`[
		true, false
	]`))

	// Parsing
	assert.NilError(t, err)

	// Item values
	assert.Equal(t, true, array.Get(0).Bool())
	assert.Equal(t, false, array.Get(1).Bool())

	// Array length
	assert.Equal(t, 2, array.Len())
}

func TestPath(t *testing.T) {
	node, err := New([]byte(`{
		"foo": {
			"bar": [
				{ "x": "a" },
				{ "x": "b" }
			]
		}
	}`))

	// Parsing
	assert.NilError(t, err)

	// Path works with object and array
	assert.Equal(t, "a", node.Get("foo", "bar", 0, "x").Str())
	assert.Equal(t, "b", node.Get("foo", "bar", 1, "x").Str())

	// Invalid path returns error node
	invalidPath := node.Get("where", "is", "this")
	assert.Equal(t, Error, invalidPath.Type())
	assert.Error(t, invalidPath.Err(), "path element not found: 'where'")

	// Passing a string key on an array node returns error node
	invalidStringKey := node.Get("foo", "bar", "buzz")
	assert.Equal(t, Error, invalidStringKey.Type())
	assert.Error(t, invalidStringKey.Err(), "string key used on non-object of type: []interface {}, key: 'buzz'")

	// Passing an integer key on an object node returns error node
	invalidIntegerKey := node.Get("foo", 0)
	assert.Equal(t, Error, invalidIntegerKey.Type())
	assert.Error(t, invalidIntegerKey.Err(), "integer key used on non-array of type: map[string]interface {}, index: 0")

	// Index out of bounds for array
	indexOutOfBounds := node.Get("foo", "bar", -1, "x")
	assert.Equal(t, Error, indexOutOfBounds.Type())
	assert.Error(t, indexOutOfBounds.Err(), "array index put of bounds: index -1 of length: 2")
	indexOutOfBounds = node.Get("foo", "bar", 2, "x")
	assert.Equal(t, Error, indexOutOfBounds.Type())
	assert.Error(t, indexOutOfBounds.Err(), "array index put of bounds: index 2 of length: 2")

}

func TestNull(t *testing.T) {
	json, err := New([]byte(`{ "foo": null }`))

	// Parsing
	assert.NilError(t, err)

	// Get the null node
	nullNode := json.Get("foo")

	// Check null value
	assert.Equal(t, nil, nullNode.Raw())

	// Check null type
	assert.Equal(t, Null, nullNode.Type())

	// Null defaults to empty string
	assert.Equal(t, "", nullNode.Str())

	// Null defaults to number zero
	assert.Equal(t, 0.0, nullNode.Num())

	// Null defaults to boolean false
	assert.Equal(t, false, nullNode.Bool())

	// Accessing a path under null value returns an error node
	assert.Equal(t, Error, nullNode.Get("bar").Type())
	assert.Error(t, nullNode.Get("bar").Err(), "string key used on non-object of type: <nil>, key: 'bar'")
}
