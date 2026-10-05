package main

import "fmt"

func main() {
	pointer()
}


func pointer() {
	i, j := 42, 2701
	var nullPointer *int
	var nullPointer2 *int

	p := &i
	fmt.Println(nullPointer)
	fmt.Println(&nullPointer2)
	fmt.Println("P pointer", p)
	fmt.Println("P value", *p)
	*p = 21
	fmt.Println("New i", i)  // see the new value of i

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer
	fmt.Println("New j", j) // see the new value of j
}