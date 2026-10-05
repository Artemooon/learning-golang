package main

import (
    "fmt"
)

func fibonacci() func() int {
    current := 0
    next := 1

    return func() int {
        result := current
        current, next = next, current + next

        return result
    }
}

func main() {
    f := fibonacci()
    for i := 0; i < 10; i++ {
        fmt.Println(f())
    }

}