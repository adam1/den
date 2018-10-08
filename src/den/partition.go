// Copyright 2018 Adam Marks

package den

import (
	"fmt"
	"sort"
	"time"
)

type Partition []int // xxx change this from int to uint8; add assertion in AllPartitions to verify n <= 255.  make an alias for uint8 to PartitionInt.


func (p Partition) String() string {
	var s string = "["
	for i, part := range p {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%d", part)
	}
	return s + "]"
}

func (p Partition) Check(degree int) bool {
	if p.Sum() == degree {
		return true
	}
	return false
}

func (p Partition) Sum() int {
	var sum int
	for _, p := range p {
		sum += p
	}
	return sum
}

func (p Partition) CycleTypeOld() CycleType {
	sum := p.Sum()
	t := make([]int, sum, sum)
	p.CycleType(t)
	return t
}

func (p Partition) CycleType(t CycleType) {
	for i := range t {
		t[i] = 0
	}
	for _, p := range p {
		t[p-1]++
	}
}

func (p Partition) Equal(q Partition) bool {
	if len(p) == len(q) {
		for i := range p {
			if p[i] != q[i] {
				return false
			}
		}
		return true
	}
	return false
}

// AllPartitions spawns a goroutine that generates the partitions and
// emits them on the returned channel as they are generated.  when all
// partitions have been generated, the goroutine ends and the channel
// is closed.
func YieldAllPartitions(degree int) chan Partition {
	yield := make(chan Partition, 100)
	go func() {
		ruleAsc(degree, yield)
		close(yield)
	}()
	return yield
}

func AllPartitions(degree int) []Partition {
	count := CountAllPartitions(degree)
	partitions := make([]Partition, count)
	var yield chan Partition = YieldAllPartitions(degree)
	i := 0
	for p := range yield {
		partitions[i] = p
		i++
	}
	return partitions
}

func CountAllPartitions(degree int) int {
	var yield chan Partition = YieldAllPartitions(degree)
	i := 0
	for range yield {
		i++
	}
	return i
}

// based on ruleAsc by jerome kelleher
//   http://homepages.ed.ac.uk/jkellehe/partitions.php
func ruleAsc(n int, yield chan Partition) {
	debug := false
	animate := false
	var a []int = make([]int, n+1) // xxx +1 not needed here?
	var k int = 1
	a[0] = 0
	a[1] = n
	for k != 0 {
		x := a[k-1] + 1
		y := a[k] - 1
		k -= 1
		for x <= y {
			a[k] = x
			y -= x
			k += 1
		}
		a[k] = x + y
		var b []int = make([]int, k+1)
		copy(b, a[:k+1])
		yield <- b
		if debug {
			fmt.Printf("k=%d x=%d y=%d a=%v b=%v\n", k, x, y, a, b)
		} else if animate {
			fmt.Printf("\033[2Ka=%v\r", a)
			time.Sleep(100*time.Millisecond)
		}
	}
}

type SortablePartitions []Partition

func (partitions SortablePartitions) Len() int {
	return len(partitions)
}

// This ordering operator is compatible with ruleAsc, in that if
// ruleAsc produces A before B, then Less(A, B) is true.  Furthermore,
// it applies across differing values of N.

func (p Partition) Less(q Partition) bool {
	for k, x := range p {
		y := 0 // imagine a suffix of infinite zeroes 
		if k < len(q) {
			y = q[k]
		}
		if x < y { // xxx hot spot
			return true
		}
		if y < x {
			return false
		}
	}
	return len(p) < len(q)
}

func (partitions SortablePartitions) Less(i, j int) bool {
	a := partitions[i]
	b := partitions[j]
	return a.Less(b)
}

func (partitions SortablePartitions) Swap(i, j int) {
	t := partitions[i]
	partitions[i] = partitions[j]
	partitions[j] = t
}

func (partitions SortablePartitions) Equal(Q SortablePartitions) bool {
	if len(partitions) == len(Q) {
		for i, p := range partitions {
			if !p.Equal(Q[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (partitions SortablePartitions) Search(p Partition) (index int, found bool) {
	index = sort.Search(len(partitions), func(i int) bool {
		val := false
		q := partitions[i] // xxx hot spot
		if p.Equal(q) || p.Less(q) {
			val = true
		}
		//log.Printf("tested i=%d L=%d p.Less(q) p=%v q=%v got=%v", i, len(partitions), p, partitions[i], val)
		return val
	})
	if index < len(partitions) && partitions[index].Equal(p) {
		found = true
	}
	return index, found
}

func (partitions SortablePartitions) String() string {
	s := ""
	for _, p := range partitions {
		s += fmt.Sprintf("%v\n", p)
	}
	return s
}
