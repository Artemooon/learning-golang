package main

import (
    "fmt"
    "math"
)

type Vertex struct {
    X, Y float64
}

// func (v Vertex) Abs() float64 {
//     return math.Sqrt(v.X*v.X + v.Y*v.Y)
// }

func (v *Vertex) Abs() float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y)
}


func (v *Vertex) Scale(f float64) {
    fmt.Println("V: ", *v)
    v.X = v.X * f
    v.Y = v.Y * f
    fmt.Println("V: ", *v)
}

func AbsFunc(v Vertex) float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// func AbsFunc(v *Vertex) float64 {
//     return math.Sqrt(v.X*v.X + v.Y*v.Y)
// }

func ScaleFunc(v *Vertex, f float64) {
    v.X = v.X * f
    v.Y = v.Y * f
}

func main() {
    v := Vertex{79.34, 98.43}
    fmt.Println("Original: ", v.Abs())
    (&v).Scale(2)
    fmt.Println("Scaled: ", v.Abs())

    ScaleFunc(&v, 2)
    fmt.Println("AbsFunc: ", AbsFunc(v))

    p := &Vertex{4, 3}
    p.Scale(3)
    ScaleFunc(p, 2)
    fmt.Println("AbsFunc: ", AbsFunc(*p))
    fmt.Println("v: ", &v)
}