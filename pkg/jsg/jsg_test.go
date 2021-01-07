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
