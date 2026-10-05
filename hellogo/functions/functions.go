package main

import (
	"fmt"
	"math"
)

func compute(fn func(float64, float64) float64, loggerFn ...func(string)) float64 {
    if len(loggerFn) > 0 && loggerFn[0] != nil {
        loggerFn[0]("Computing result for: ")
    }
	return fn(3, 4)
}

func main() {
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}

    logger := func(text string) {
        fmt.Println("Logged: ", text)
    }

	fmt.Println(hypot(5, 12))

	fmt.Println(compute(hypot, logger))
	fmt.Println(compute(math.Pow, logger))
}


func main1() {



}