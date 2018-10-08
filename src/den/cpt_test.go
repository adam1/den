// Copyright 2018 Adam Marks

package den

import (
	"log"
	"math/rand"
	"testing"
)

// xxx todo: refactor to add support for explicity, nonrandom, test
// cases, and a case that seems broken in check-pre-extensions: d=2 u=(2)
func TestLogarithms(t *testing.T) {
	debug := false
	maxPower := 1000
	maxDegree := 30
	iterations := 100
	if testing.Short() {
		maxDegree = 22
		iterations = 13
	}
	for i := 0; i < iterations; i++ {
		degree := 1 + rand.Intn(maxDegree)
		u := RandomCycleType(degree)
		k := rand.Intn(maxPower)
		v := u.PowerOld(k)
		o := u.Order()
		r := k % o
		if r == 0 {
			// We don't include the zeroth powers in
			// logarithm results, hence kick it up to the
			// first equivalent positive power.
			r = o
		}
		if debug {
			log.Printf("i=%d d=%d u=%v k=%v o=%d r=%d v=%v", i, degree, u, k, o, r, v)
		}
// 		v := u.PowerOld(k)
// 		reducedPower := k % u.Order()
// 		log.Printf("xxx i=%d d=%d u=%v k=%v r=%d v=%v", i, degree, u, k, reducedPower, v)
// 		//log.Printf("xxx pft:\n%v", P)
		C := New_CPT(degree)
		C.Generate()
		logarithms := C.Logarithms(v)
		found := false
		for _, x := range logarithms {
			if x.Base.Equal(u) && x.Power == r {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Logarithm not found; i=%d d=%d u=%v k=%v o=%d r=%d logarithms=%v", i, degree, u, k, o, r, logarithms)
		}
	}
}
