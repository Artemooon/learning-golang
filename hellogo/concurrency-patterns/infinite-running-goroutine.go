package main

import (
    "fmt"
    "time"
)

func main() {

    go func() {
        for {
            select {
                default:
                    fmt.Println("goroutine working")
            }
        }
    }()

    time.Sleep(time.Second * 2)
}