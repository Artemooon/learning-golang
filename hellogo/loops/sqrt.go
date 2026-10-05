package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := float64(1)
	for {
		prev := z
		z -= (z*z - x) / (2.0*z)
		if math.Abs(z-prev) < 0.0000001 {
			return z
		}
		fmt.Println("Z: ", z)
	}
	return z
}

func main() {
	fmt.Println(Sqrt(4040304004304343))
}
