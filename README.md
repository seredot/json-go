# JSON-GO  ********* WIP *********

Go Library for reading, writing, modifying JSON like data.

```go
json := byte[](`{
    "name": "John Doe",
    "age": 37,
    "toys": null,
    "children": [
        {"name": "Irene"},
        {"name": "Alan"}
    ]
}`)

// Get the root node:
john := jsg.New(json)

// Get values from fields
name := john.Get("name").Str() // "John Doe"
age := john.Get("age").Int() // 37
daughtersName := john.Get("children").Get(0).Get("name").Str() // "Irene"

// Or use paths
daughtersName := john.Get("children", 0, "name").Str() // "Irene"

// Get an array
children := john.Get("children")

// Invalid paths return empty values safely
sonsName := children.Get(5).Get("name").Str() // ""
sonsName = children.Get(5, "name").Str() // ""

// Check for errors
err := children.Get(5).Get("name").Err() // "array index out of bounds"
err = children.Get(5, "name").Err() // "array index out of bounds"

// Check for types
isArray := children.Type() == jsg.Array // true
isNumber := john.Get("age").Type() == jsg.Number // true
isNull := john.Get("toys").Type() == jsg.Null // true
isNull = john.Get("toys").Raw() == nil // true

// Set a primitive field
john.Set("age", 38)
john.Set("spouse", "Jane Doe")

// Create object and set field it's fields
parents := john.Set("parents", jsg.NewObject())
parents.Set("father", "Mark")
parents.Set("mother", "Rosetta")
parents.Set("married", true)

// Add an item to an array
baby := children.Push(jsg.NewObject())
baby.Set("name", "Ada")

// Get the length of an array
length := children.Len()

// Serialize back to JSON string
output := string(dt.SerializeIndent("\t"))
```

License: MIT
