package main

import (
    "fmt"
    "math"
)


type Vertex struct {
    Long float64
    Lat float64
}

// func (v Vertex) Abs() float64 {
//     return math.Sqrt(v.Lat*v.Long)
// }

func Abs(v Vertex) float64 {
    return math.Sqrt(v.Lat*v.Long)
}

func main() {
    v := Vertex{12.0, 14.4}
    fmt.Println("Abs", Abs(v))
}