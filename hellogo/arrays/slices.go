package main

import "fmt"

func main() {
	primes := [6]int{2, 3, 5, 7, 11, 13}

	var s []int = primes[1:6]
	fmt.Println(s)

	names := [4]string{
		"John",
		"Paul",
		"George",
		"Ringo",
	}
	fmt.Println(names)

	a := names[0:2]
	b := names[1:3]
	fmt.Println(a, b)

	b[0] = "XXX"
	fmt.Println(a, b)
	fmt.Println(names)

	sliceLiterals()
	sliceSlicing()
}


func sliceLiterals() {
	q := []int{1,2,3,4,5}
	fmt.Println(q)

	r := []bool{true, false, true, true, false, true}
	fmt.Println(r)

	s := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	fmt.Println(s)

	sliceStrings := []string{"Hi", "it's", "me"}
	fmt.Println(sliceStrings)
}


func sliceSlicing() {
	arr := [4]int{1,2,3,4}
	s := []int{2, 3, 5, 7, 11, 13}
	b := []bool{true, false, true}


	fmt.Println(s)
	s = s[1:4]
	fmt.Println(s)

	s = s[:2]
	fmt.Println(s)

	s = s[1:]
	fmt.Println(s)

	b = b[:]
	fmt.Println(b)

	b = b[0:]
	fmt.Print(b)

	sliceFromArr := arr[1:4]
	fmt.Println("\n", sliceFromArr)
	sliceFromArr[0] = 3
	fmt.Println(sliceFromArr)
	fmt.Println(arr)
	fmt.Println(sliceFromArr, cap(sliceFromArr))
	sliceFromArr = sliceFromArr[:cap(sliceFromArr)]
	fmt.Println(sliceFromArr, cap(sliceFromArr))

	t := make([]byte, len(s), (cap(s)+1)*2) // +1 in case cap(s) == 0
	for i := range s {
			t[i] = s[i]
	}
	s = t
}