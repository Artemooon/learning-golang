// package main
//
// import (
//     "fmt"
// )
//
// type I interface {
//     M()
// }
//
// type Plusable interface {
//     Plus2()
// }
//
// type T struct {
//     S string
// }
//
// func (t T) M() {
//     fmt.Println(t.S)
// }
//
// type F float64
//
// func (f *F) Plus2() {
//     *f = *f + 2
// }
//
// func main() {
//     var i I = T{"hello"}
//     value := F(2.40)
//
//     var plusable Plusable = &value
//     i.M()
//     plusable.Plus2()
//     fmt.Println("plusable", *plusable.(*F))
// }


package main

import (
	"fmt"
	"math"
)

type I interface {
	M()
}

type T struct {
	S string
}

func (t *T) M() {
    if (t == nil) {
        fmt.Println("<Nil>")
        return
    }
	fmt.Println(t.S)
}

type F float64

func (f F) M() {
	fmt.Println(f)
}

func main() {
	var i I

    describe(i)
    fmt.Println(i == nil)
    var t *T
	i = t
	describe(i)
	iv := i
    fmt.Println(iv[0] == nil)
    t, ok := i.(*T)
	i.M()

	i = F(math.Pi)
	describe(i)
	i.M()
}

func describe(i I) {
	fmt.Printf("(%v, %T)\n", i, i)
}
