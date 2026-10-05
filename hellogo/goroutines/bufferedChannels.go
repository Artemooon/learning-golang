package main

import (
    "fmt"
    "time"
)


func removeDuplicates(duplicatesSlice []int, ch chan []int) {
    duplicatesMap := make(map[int]bool)
    unique := make([]int, 0, len(duplicatesSlice))

    for _, v := range duplicatesSlice {
        if duplicatesMap[v] {
            continue
        }
        duplicatesMap[v] = true
        unique = append(unique, v)
    }

    time.Sleep(2 * time.Second)

    fmt.Println("worker", id, "attempting to send")
    ch <- unique
    fmt.Println("worker", id, "sent result")
}

func main() {
    ch := make(chan []int, 2)
    mySlice := []int{1,2,3,2,0,0,9,5,20,20,143,52,21,134,143,12,21,9,4,9,2,2,7}

    fmt.Println("mySlice", mySlice)
    for i := 0; i < 5; i++ {
        go removeDuplicates(mySlice[i:i+10], ch)
    }
    time.Sleep(3 * time.Second)

    fmt.Println("queued:", len(ch)) // normally 2
    fmt.Println("capacity:", cap(ch)) // 2

}