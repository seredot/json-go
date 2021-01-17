package jsg

import (
	"fmt"
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
	assert.Error(t, indexOutOfBounds.Err(), "array index out of bounds: index -1 of length: 2")
	indexOutOfBounds = node.Get("foo", "bar", 2, "x")
	assert.Equal(t, Error, indexOutOfBounds.Type())
	assert.Error(t, indexOutOfBounds.Err(), "array index out of bounds: index 2 of length: 2")
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
	assert.ErrorContains(t, nullNode.Get("bar").Err(), "string key used on non-object of type: ")
}

func TestUndefined(t *testing.T) {
	json, err := New([]byte(`{ "foo": 1 }`))

	// Parsing
	assert.NilError(t, err)

	// Get the undefined node
	nullNode := json.Get("bar")

	// Undefined field returns error node
	assert.Equal(t, Error, nullNode.Type())

	// Check the error message
	assert.ErrorContains(t, nullNode.Err(), "path element not found: ")
}

func TestError(t *testing.T) {
	json, err := New([]byte(`{}`))

	// Parsing
	assert.NilError(t, err)

	// Get the error node
	errorNode := json.Get("foo")

	// Undefined field returns error node
	assert.Equal(t, Error, errorNode.Type())

	// Error defaults to empty string
	assert.Equal(t, "", errorNode.Str())

	// Error defaults to number zero
	assert.Equal(t, 0.0, errorNode.Num())

	// Error defaults to boolean false
	assert.Equal(t, false, errorNode.Bool())

	// Error message
	assert.ErrorContains(t, errorNode.Err(), "path element not found: ")
	assert.Error(t, fmt.Errorf("path element not found: 'foo'"), errorNode.Raw().(error).Error())

	// Set using an invalid type (int here)
	err = NewObject().Set("foo", 1).Err()
	assert.Error(t, err, errorInvalidType)
}

func TestSet(t *testing.T) {
	// Create a new object
	root := NewObject()

	// Use the setter
	root.Set("foo", "bar")

	// Read back the value
	assert.Equal(t, "bar", root.Get("foo").Str())

	// Create a new object
	array := NewArray()

	// Use the setter
	array.Set(2, "fun")
	array.Set(1, "is")
	array.Set(0, "go")

	// Read back values
	assert.Equal(t, "go", array.Get(0).Str())
	assert.Equal(t, "is", array.Get(1).Str())
	assert.Equal(t, "fun", array.Get(2).Str())

	// Set object field with an array index key
	err := NewObject().Set(0, "foo").Err()
	assert.ErrorContains(t, err, "integer key used on non-array of type:")

	// Set array item with a string key
	err = NewArray().Set("foo", "bar").Err()
	assert.ErrorContains(t, err, "string key used on non-object of type:")

	// Set array item with a negative index
	err = NewArray().Set(-1, "bar").Err()
	assert.ErrorContains(t, err, "array index out of bounds: index -1")

	// Set using an invalid key type (other than string or int)
	err = NewObject().Set(false, "bar").Err()
	assert.ErrorContains(t, err, "path params must be string or integer")
}

func TestDel(t *testing.T) {
	object := NewObject()
	object.Set("foo", "bar")

	// Delete key from object
	err := object.Del("foo")
	assert.NilError(t, err)

	// Deleted field should return error node
	assert.Equal(t, Error, object.Get("foo").Type())

	array := NewArray()
	array.Set(0, 1.0)
	array.Set(1, 2.0)
	array.Set(2, 3.0)

	// Delete first item from array
	err = array.Del(0)
	assert.NilError(t, err)

	// Delete the last item
	err = array.Del(array.Len() - 1)
	assert.NilError(t, err)

	// There should be 1 item left
	assert.Equal(t, 1, array.Len())
	assert.Equal(t, 2.0, array.Get(0).Num())

	// Deleted negative index return error node
	assert.ErrorContains(t, array.Del(-1), "array index out of bounds: index -1")

	// Deleted out of bounds index return error node
	assert.ErrorContains(t, array.Del(10), "array index out of bounds: index 10")
}
