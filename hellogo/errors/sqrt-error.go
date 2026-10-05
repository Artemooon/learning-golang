package main


import (
    "fmt"
    "math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
    return fmt.Sprintf("cannot Sqrt negative number: %f", float64(e))
}

func Sqrt(x float64) (float64, error) {
    if x <= 0 {
        return 0.0, ErrNegativeSqrt(x)
    }
	z := float64(1)
	for {
		prev := z
		z -= (z*z - x) / (2.0*z)
		if math.Abs(z-prev) < 0.0000001 {
			return z, nil
		}
		fmt.Println("Z: ", z)
	}
	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
    fmt.Println(Sqrt(0))
}