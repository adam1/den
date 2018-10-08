// Copyright 2018 Adam Marks

package den

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"strconv"
)

// sagan format [1^m_1, 2^m_2, ..., n^m_n] where m_k is the number of
// cycles of length k
type CycleType []int

func NewCycleType(sagan []int) *CycleType {
	x := make([]int, len(sagan))
	copy(x, sagan)
	return (*CycleType)(&x)
}

func RandomCycleType(degree int) *CycleType {
	// at each step, pick next random cycle length <= remaining
	// then pick a random quantity of them, up to remaining/length
	result := make([]int, degree)
	remaining := degree
	for remaining > 0 {
		L := 1 + rand.Intn(remaining)
		q := 1 + rand.Intn(int(remaining/L))
		result[L-1] += q
		remaining -= L*q
	}
	return NewCycleType(result)
}

func (lambda *CycleType) Set(value string) error {
	*lambda = (*lambda)[0:0]
	for _, s := range strings.Split(value, ",") {
		j, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*lambda = append(*lambda, j)
	}
	return nil
}

func (ct *CycleType) Copy() *CycleType {
	return NewCycleType(*ct)
}

func (ct *CycleType) Equal(x *CycleType) bool {
	if ct == nil && x == nil { 
		return true; 
	}
	if ct == nil || x == nil { 
		return false; 
	}
	if len(*ct) != len(*x) {
		return false
	}

	for i := range *ct {
		if (*ct)[i] != (*x)[i] {
			return false
		}
	}
	return true
}

func (ct *CycleType) DegreeOld() int {
	d := 0
	for i, m := range *ct {
		d += (i+1) * m
	}
	return d
}

func (t *CycleType) Degree() int {
	return len(*t)
}

func (lambda *CycleType) Order() int {
	present := make([]int, 0)
	for i, x := range *lambda {
		if x > 0 {
			present = append(present, i+1)
		}
	}
	return LCM(present)
}

func (t *CycleType) Orderb(order *big.Int) {
	present := make([]int, 0)
	for i, x := range *t {
		if x > 0 {
			present = append(present, i + 1)
		}
	}
	LCMb(present, order)
}

func (t *CycleType) PowerOld(k int) *CycleType {
	u := make([]int, len(*t))
	t.Power(k, u)
	return (*CycleType)(&u)
}

func (t *CycleType) Power(k int, u CycleType) {
	for i := range u {
		u[i] = 0
	}
	for i := 1; i <= len(*t); i++ {
		// -1's are for 0-based indexing into the int slice t
		m := (*t)[i - 1]
		if m > 0 {
			f := GCD(i, k)
			u[i/f - 1] += f * m
		}
	}
}

func (t *CycleType) IsIdentity() bool {
	for i, c := range *t {
		if i == 0 {
			if c == 0 {
				return false
			}
		} else {
			if c != 0 {
				return false
			}
		}
	}
	return true
}

func (t *CycleType) Partition(p *Partition, buf []int) {
	k := 0
	for i, m := range *t {
		for j := 0; j < m; j++ {
			buf[k] = i+1
			k++
		}
	}
	*p = buf[:k]
}

func (t *CycleType) CardinalityOfConjugacyClass() *big.Int {
	var k *big.Int = Factorial(t.Degree())
	k.Div(k, t.CardinalityOfCentralizer())
	return k
}

func (t *CycleType) CardinalityOfCentralizer() *big.Int {
	k := big.NewInt(1)
	for i := 0; i < t.Degree(); i++ {
		s := i + 1
		m := (*t)[i]
		k.Mul(k, Exp(s, m))
		k.Mul(k, Factorial(m))
	}
	return k
}

func (ct *CycleType) HashKeyString() string {
	return fmt.Sprint(*ct)
}

func (ct *CycleType) String() string {
	return ct.StringWithCarets()
}

func (ct *CycleType) StringWithTilde() (s string) {
	return ct.StringWithQuantityMarker("~")
}

func (ct *CycleType) StringWithColons() (s string) {
	return ct.StringWithQuantityMarker(":")
}

func (ct *CycleType) StringWithCarets() (s string) {
	return ct.StringWithQuantityMarker("^")
}

func (ct *CycleType) StringWithQuantityMarker(marker string) (s string) {
	for i := len(*ct) - 1; i >= 0; i-- {
		k := i + 1
		m := (*ct)[i]
		if m > 0 {
			if len(s) == 0 {
				s = "("
			} else {
				s += ","
			}
			s += fmt.Sprintf("%d", k)
			if m > 1 {
				s += fmt.Sprintf("%s%d", marker, m)
			}
		}
	}
	return s + ")"
}

func (ct *CycleType) StringForPartitionWithoutOneCycles() (s string) {
	for i := len(*ct) - 1; i >= 1; i-- {
		k := i + 1
		n := (*ct)[i]
		if n > 0 && len(s) > 0 {
			s += ","
		}
		for j := 0; j < n; j++ {
			if j > 0 {
				s += ","
			}
			s += fmt.Sprintf("%d", k)
		}
	}
	if s == "" {
		s = "1"
	}
	return s
}


// xxx slightly weird; maybe wrap the slice, degree
func (lambda *CycleType) Pad(degree int) CycleType {
	padded := make([]int, degree)
	copy(padded, *lambda)
	return padded
}

func (ct *CycleType) PreExtensions() []*CycleType {
	result := make([]*CycleType, 0)
	for i, m := range *ct {
		if m > 0 {
			preExtension := ct.Copy()
			(*preExtension)[i] -= 1
			if i > 0 {
				(*preExtension)[i-1] += 1
			}
			result = append(result, preExtension)
		}
	}
	return result
}

// terminology: the height of a type is the number of powers of the
// type that equal the type itself.
func (t *CycleType) markMethodHeight(height *big.Int) {
	height.SetInt64(1)
	bigOne := big.NewInt(1)
	k := 2
	var u CycleType = make([]int, t.Degree())
	for {
		t.Power(k, u)
		//log.Printf("%v^%v = %v", t, k, &u)
		if u.IsIdentity() {
			break
		}
		if u.Equal(t) {
			height.Add(height, bigOne)
		}
		k++
	}
	//log.Printf("t=%v k=%v mark height=%v", t, k, height)
}

func (t *CycleType) totientMethodHeight(height *big.Int) {
	order := big.NewInt(0)
	t.Orderb(order)
	Totient(order, height)
	//log.Printf("t=%v order=%v totient height=%v", t, order, height)
}


