package main

import (
	"fmt"
	"runtime"
	"time"
	"math"
)

func main() {
	fmt.Print("Go runs on ")
	fmt.Println("Current OS: ", runtime.GOOS)
	switch os := ""; os {
	case "darwin":
		fmt.Println("macOS.")
	case "linux":
		fmt.Println("Linux.")
	case "windows":
		fmt.Println("GG WP.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Println("No name")
		fmt.Printf("%s.\n", os)
	}
	main2()
	main3()
	mainDefer()
	// mainCount()
	var test = c()
	fmt.Println("test", test)
	pointer()
}

func main3() {
	fmt.Println("When's Monday?")
	today := time.Now().Weekday()
	fmt.Println("Sunday", time.Sunday)
	fmt.Println("today + 0", today + 2)
	switch time.Monday {
	case today + 1:
		fmt.Println("Today.")
	case today + 2:
		fmt.Println("Tomorrow.")
	case today + 3:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}
}


func main2() {
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}


func mainDefer() int {
	defer fmt.Println("world")
	defer fmt.Println(math.Round(0.94))

	fmt.Println("hello")
	return 0;
}

func mainCount() {
	fmt.Println("counting")

	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}

	fmt.Println("done")
}

func c() (i int) {
    defer func() { i++ }()
    return 1
}
