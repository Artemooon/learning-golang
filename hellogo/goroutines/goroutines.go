package main

import (
    "fmt"
    "time"
//     "runtime"
)

func say(s string) {
    for i := 0; i < 5; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Println(s)
    }
}

func blocking1() {
    fmt.Println("Blocking 1 started")
    time.Sleep(time.Second * 1)
    fmt.Println("Blocking 1 ended")
}

func blocking2(c chan int) {
    fmt.Println("Blocking 2 started")
    time.Sleep(time.Second * 3)
    fmt.Println("Blocking 2 ended")
    a := 5
    c <- a
}

func blocking3(c chan int) {
    fmt.Println("Blocking 3 started")
    time.Sleep(time.Second * 3)
    fmt.Println("Blocking 3 ended")
    a := 10
    c <- a
}

func sum(a int, b int) {
    fmt.Println("sum", a + b)
}

func main() {

        go sum(1, 2)
        go sum(3, 4)

        time.Sleep(time.Second * 2)

        fmt.Println("Main function ended")
//     a := 1
//     myChan := make(chan int)
//     myChan2 := make(chan int)
//
//     blocking1()
//     go blocking2(myChan)
//     go blocking3(myChan2)
// //     fmt.Println(a)
//     a2 := <- myChan
//     a3 := <- myChan2
//     fmt.Println("a2", a2)
//     fmt.Println("a3", a3)
//     fmt.Println(runtime.GOMAXPROCS(0))
//     go say("world")
//     go say("hello")
}