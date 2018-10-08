// Copyright 2018 Adam Marks

package den

import (
	"fmt"
)

// partition fold table (PFT) xxx factor out / merge into
// cycle_type.go... PFT object may remain as a simple wrapper in cases
// where the whole table is desired.

type PFT struct {
	degree int
	lambda CycleType
	data []CycleType
}

func NewPFT(degree int, lambda CycleType) *PFT {

	var X *PFT = new(PFT)

	X.degree = degree
	X.lambda = lambda.Pad(degree)

	return X
}

// xxx factor out CycleType.Power
func (X *PFT) Generate() {
	order := X.lambda.Order()
	//fmt.Printf("generating PFT degree=%v lambda=%v order=%v\n", X.degree, X.lambda, order)
	X.data = make([]CycleType, order)
	// note annoying conversion from mathematical 1-based notation
	// to programmatical 0-based mechanics.
	for a := 0; a < order; a++ {
		X.data[a] = make([]int, X.degree)
		// xxx this copy is redundant; handled by the third case below?
		if a == 0 {
			copy(X.data[0], X.lambda)
		} 
		for b := 0; b < X.degree; b++ {
			// for the sake of my sanity
			var i = a + 1
			var j = b + 1
			k := GCD(i, j)
			if j > 1 && k == j {
				// absorb cycle
				X.data[a][0] += j * X.lambda[b]
			} else if k > 1 {
				// split cycle
				X.data[a][j/k-1] += k * X.lambda[b]
			} else {
				X.data[a][b] = X.lambda[b]
			}
		}
	}
}

func (X *PFT) Power(k int) *CycleType {
	return &X.data[(k-1) % len(X.data)]
}

// xxx redundant with CycleType.Diameter
// func (X *PFT) Order() int {
// 	return len(X.data)
// }

func (X *PFT) String() string {
	var s string
	// header
	s += "         "
	for b := 0; b < X.degree; b++ {
		s += fmt.Sprintf("%2d  ", b+1)
	}
	s += "\n"
	s += "      +--"
	for b := 0; b < X.degree; b++ {
		s += "----"
	}
	s += "\n"
	for a, _ := range X.data {
		for b, z := range X.data[a] {
			if b == 0 {
				// sider
				s += fmt.Sprintf("%6d|  ", a+1)
			} else {
				s += "  "
			}
			if z > 0 {
				s += fmt.Sprintf("%2d", z)
			} else {
				s += "  "
			}
		}
		s += "\n"
	}
	return s
}

func (X *PFT) Check() {
	// verify that each row is a partition of the degree
	for a, _ := range X.data {
		var sum int
		for b, z := range X.data[a] {
			sum += (b+1)*z
		}
		if sum != X.degree {
			panic(fmt.Sprintf("bad partition; sum=%v row=%v %v", sum, a, X.data[a]))
		}
	}
}
