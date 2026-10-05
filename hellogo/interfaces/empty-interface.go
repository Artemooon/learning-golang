package main

import (
    "fmt"
)

type I interface {}

type myStruct struct {

}

func main() {
	var i interface{}
	describe(i)

    i = 42
    describe(i)

    i = "hello"
    describe(i)

    i = false
    describe(i)
}

func describe(i interface{}) {
    fmt.Printf("Interface: (%v %T)\n", i, i)
}