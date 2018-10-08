// Copyright 2018 Adam Marks

package den

import (
	"math/big"
	"testing"
)

func TestLCM(t *testing.T) {

	tcase := func(v []int, expected int) {

		lcm := LCM(v)

		if lcm != expected {
			t.Errorf("expected=%v got=%v", expected, lcm)
		}
	}

	tcase([]int{}, 0)
	tcase([]int{0}, 0)
	tcase([]int{1,2}, 2)
	tcase([]int{2,1}, 2)
	tcase([]int{7,7,7,7,7,7,7,7}, 7)
	tcase([]int{7,7,7,7,7,7,7,3,7,7,7,7,7,7,7,7,7,7,7,7,7}, 21)
	tcase([]int{2,2,5,6}, 30)
	tcase([]int{2,4,5,6}, 60)
	tcase([]int{2,12,5,36}, 180)
}

func TestFactorial(t *testing.T) {

	tcase := func(n int, expected *big.Int) {

		f := Factorial(n)

		if f.Cmp(expected) != 0 {
			t.Errorf("Factorial(%v) expected=%v got=%v", n, expected, f)
		}
	}

	tcase(0, big.NewInt(1))
	tcase(1, big.NewInt(1))
	tcase(2, big.NewInt(2))
	tcase(3, big.NewInt(6))
	tcase(4, big.NewInt(24))
	tcase(7, big.NewInt(5040))
	tcase(13, big.NewInt(6227020800))
}

func TestExp(t *testing.T) {

	tcase := func(a, b int, expected *big.Int) {

		x := Exp(a, b)

		if x.Cmp(expected) != 0 {
			t.Errorf("Exp(%v, %v) expected=%v got=%v", a, b, expected, x)
		}
	}

	tcase(0, 0, big.NewInt(1))
	tcase(0, 1, big.NewInt(0))
	tcase(1, 0, big.NewInt(1))
	tcase(9, 0, big.NewInt(1))
	tcase(1, 1, big.NewInt(1))
	tcase(1, 2, big.NewInt(1))
	tcase(2, 1, big.NewInt(2))
	tcase(2, 2, big.NewInt(4))
	tcase(2, 3, big.NewInt(8))
	tcase(3, 3, big.NewInt(27))
	// xxx add big case
}
