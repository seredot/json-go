# JSON-GO  ********* WIP *********

Go Library for reading, writing, modifying JSON like data.

```go
json := byte[](`{
    "name": "John Doe",
    "age": 37,
    "children": [
        {"name": "Irene"},
        {"name": "Alan"}
    ]
}`)

// Get the root node:
john := jsg.New(json)

// Get values from fields
name := john.Get("name").String() // "John Doe"
age := john.Get("age").Int()      // 38 
daughtersName := john.Get("children").At(0).Get("name").String() // "Irene"

// Or use paths
daughtersName := john.Get("children", 0, "name").String() // "Irene"

// Get an array
children := john.get("children")

// Invalid paths return empty values safely
sonsName := children.At(5).Get("name").String() // ""
sonsName := children.Get(5, "name").String() // ""

// Check for errors
err := children.At(5).Get("name").Err() // "index out of bounds"
err := children.Get(5, "name").Err() // "index out of bounds"

// Check for types
isArray := children.Type() == jsg.Array
isNumber := john.Get("age").Type() == jsg.Number

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
