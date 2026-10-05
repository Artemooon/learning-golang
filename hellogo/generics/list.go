package main

import "fmt"

// List represents a singly-linked list that holds
// values of any type.
type List[T any] struct {
	next *List[T]
	val  T
}

func main() {
    node0 := List[int]{val: 0}
    node1 := List[int]{&node0, 1}

    fmt.Println(node1.val)
}
