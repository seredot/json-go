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
	assert.Equal(t, "foo", object.Get("string").String())

	// Number value
	assert.Equal(t, 123.4, object.Get("number").Number())

	// Boolean value
	assert.Equal(t, true, object.Get("boolean").Boolean())
}
