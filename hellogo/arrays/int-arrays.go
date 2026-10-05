package main

import "fmt"

func main() {
	var stringsArr [5]string;
	stringsArr[0] = "a"
	stringsArr[2] = "b"
	stringsArr[4] = "d"
	fmt.Println("Strings array: ", stringsArr)


	intsArr := [3]int{1,2,3}
	fmt.Println("Integers array:", intsArr)

	var emptyArray [10]int;
	fmt.Println(emptyArray)

	var dynamicLengthArray = [...]string{"Artem", "Jhon"}
	fmt.Println(dynamicLengthArray)
}