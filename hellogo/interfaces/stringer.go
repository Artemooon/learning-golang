package main

import "fmt"

type Person struct {
    Name string
    Age int
    Gender string
}

func (p Person) String() string {
    return fmt.Sprintf("%+v %+v (%v years)", p.Name, p.Gender, p.Age)
}

func main() {
    p1 := Person{"Artem", 29, "M"}
    p2 := Person{"Lana", 34, "F"}
    fmt.Println(p1, p2)
}