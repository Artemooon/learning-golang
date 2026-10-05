package main

import (
    "golang.org/x/tour/tree"
    "fmt"
)

type Tree struct {
    Left  *Tree
    Value int
    Right *Tree
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
    if t == nil {
        return
    }

    Walk(t.Left, ch)
    ch <- t.Value
//     fmt.Println(t.Value)
    Walk(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        Walk(t1, ch1)
        close(ch1)
    }()

    go func() {
        Walk(t2, ch2)
        close(ch2)
    }()

    for {
        value1, ok1 := <-ch1
        value2, ok2 := <-ch2

        if ok1 == false && ok2 == false {
            return true
        }

        if ok1 != ok2 {
            return false // one tree is smaller than another
        }

        if value1 != value2 {
            return false // value does not match
        }
    }

}

func main() {
    ch := make(chan int)

    go func() {
        Walk(tree.New(22), ch)
        close(ch)
    }()

    fmt.Println("Is same?", Same(tree.New(1), tree.New(1)))
    fmt.Println("Is same?", Same(tree.New(2), tree.New(4)))

    for value := range ch {
        fmt.Println(value)
    }
}
