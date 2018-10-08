// Copyright 2018 Adam Marks

package den

import (
	"math"
	"math/big"
	"math/rand"
)

func ShuffleArray(v []int) []int {
	L := len(v)
	x := make([]int, L)
	copy(x, v)
	transpositions := 2*L
	for i := 0; i < transpositions; i++ {
		a := rand.Intn(L)
		b := rand.Intn(L)
		s := x[a]
		x[a] = x[b]
		x[b] = s
	}
	return x
}

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func LCM(v []int) int {
	if len(v) == 0 {
		return 0
	}
	z := make([]int, len(v))
	copy(z, v)
	for true {
		// determine if uniform and find min
		var min int = math.MaxInt32
		var minCol int
		var prev int
		var uniform bool = true

		for i, x := range z {
			if i > 0 && x != prev {
				uniform = false
			}
			if x < min {
				min = x
				minCol = i
			}
			prev = x
		}
		if uniform {
			break
		}
		// grow min col
		z[minCol] += v[minCol]
	}
	return z[0]
}

func LCMb(v []int, result *big.Int) {
	// future: optmize mem
	result.SetInt64(0)
	if len(v) == 0 {
		return
	}
	z := make([]*big.Int, len(v))
	for i := range v {
		z[i] = big.NewInt(int64(v[i]))
	}
	bigNegativeOne := big.NewInt(-1)
	for true {
		// determine if uniform and find min
		var min *big.Int = big.NewInt(-1)
		var minCol int
		var prev *big.Int = big.NewInt(0)
		var uniform bool = true

		for i, x := range z {
			if i > 0 && x.Cmp(prev) != 0 {
				uniform = false
			}
			if x.Cmp(min) < 0 || min.Cmp(bigNegativeOne) == 0 {
				min.Set(x)
				minCol = i
			}
			prev.Set(x)
		}
		if uniform {
			break
		}
		// grow min col
		z[minCol].Add(z[minCol], big.NewInt(int64(v[minCol])))
	}
	result.Set(z[0])
}

func Factorial(n int) *big.Int {
	if n < 0 {
		return big.NewInt(0)
	}
	if n == 0 {
		return big.NewInt(1)
	}
	f := big.NewInt(1)
	for n > 1 {
		f.Mul(f, big.NewInt(int64(n)))
		n--
	}
	return f
}

func Exp(a, b int) *big.Int {
	x := big.NewInt(0)
	return x.Exp(big.NewInt(int64(a)), big.NewInt(int64(b)), nil)
}

func Totient(a, result *big.Int) {
	result.SetInt64(1)
	bigOne := big.NewInt(1)
	g := big.NewInt(0)
	for i := big.NewInt(2); i.Cmp(a) == -1; i.Add(i, bigOne) {
		g.GCD(nil, nil, a, i)
		//log.Printf("a=%v i=%v gcd=%v", a, i, g)
		if g.Cmp(bigOne) == 0 {
			result.Add(result, bigOne)
		}
	}
}
