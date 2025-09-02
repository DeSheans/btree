# BTree Implementation in Go

A generic B-tree implementation with support for custom key types and comparison functions.

## Installation

```bash
go get github.com/DeSheans/btree
```

## Basic Usage

```go
package main

import (
    "fmt"
    "github.com/DeSheans/btree"
)

func main() {
    // Create a B-tree with order 3
    tree := btree.New[int, int](3)

    // Add key-value pairs
    tree.TryAdd(10, 100)
    tree.TryAdd(5, 50)
    tree.TryAdd(15, 150)
    tree.TryAdd(3, 30)
    tree.TryAdd(7, 70)

    // Retrieve values
    if ok, value := tree.TryGet(5); ok {
        fmt.Printf("Value for key 5: %d\n", value) // Output: 50
    }

    // Update values
    tree.TryUpdate(5, 55)

    // Remove keys
    tree.TryRemove(3)

    // Traverse in-order
    items := tree.InorderTraverse()
    for _, item := range items {
        fmt.Printf("Key: %d, Value: %d\n", item.Key, item.Value)
    }

    // Get total count
    fmt.Printf("Total items: %d\n", tree.Count())
}
```

## Custom Key Types

```go
type Person struct {
    Name string
    Age  int
}

// Create B-tree with custom struct keys
tree := btree.NewWith[Person, int](2, func(a, b Person) int {
    if a.Age < b.Age {
        return -1
    } else if a.Age > b.Age {
        return 1
    }
    return 0
})
```
