package main

import (
	"fmt"
)

func main() {
	a := make([]int, 5, 50)  // len(a)=5
	arr := [50]int{1,2,3,4,5,6,10}
	arrPointer := &arr;
	var intPointer *int = &arr[0]
	*intPointer = 55

	fmt.Println("A slice", a)
	fmt.Println("A slice len", len(a))
	fmt.Println("A slice cap", cap(a))
	fmt.Println("arr" , arr)
	fmt.Println("Arr len", len(arr))
	fmt.Println("Arr cap", cap(arr))
	fmt.Println("Arr pointer", *arrPointer)
	fmt.Println("Int pointer", *intPointer)
}