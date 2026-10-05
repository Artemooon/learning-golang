package main

import (
	"fmt"
)

type Vertex struct {
	X int
	Y int
}


func main() {
	v := Vertex{4, 5}
	v.X = 4
	fmt.Println("Vertex X", v.X)
	fmt.Println("Vertex Y", v.Y)
}