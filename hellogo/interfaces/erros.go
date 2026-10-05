package main

import (
    "fmt"
    "time"
)

type MyError struct {
    When time.Time
    What string
}

func (myErr MyError) Error() string {
    return fmt.Sprintf("at %v, %s", myErr.When, myErr.What)
}

func run() error {
    return MyError{
        time.Now(),
        "id didn't work",
    }
}

func main() {
    if err := run(); err != nil {
        fmt.Println(err)
    }
}