// Copyright 2018 Adam Marks

package den

import (
	"sort"
	"testing"
)

func TestAllPartitions(t *testing.T) {
	// http://oeis.org/A000041/b000041.txt
	knownPartitionNumbers := []int{
		1, 1, 2, 3, 5, 7, 11, 15, 22, 30, 42, 56, 77, 101, 135, 176, 231, 297, 385, 490, 627, 
		792, 1002, 1255, 1575, 1958, 2436, 3010, 3718, 4565, 5604, 6842, 8349, 10143, 12310, 14883, 17977, 21637, 
		26015, 31185, 37338, 44583, 53174, 63261, 75175, 89134, 105558, 124754, 147273, 173525, 204226, 239943, 
		281589, 329931, 386155, 451276, 526823, 614154, 715220, 831820, 966467, 1121505, 1300156, 1505499, 1741630, 
		2012558, 2323520, 2679689, 3087735, 3554345, 4087968, 4697205, 5392783, 6185689, 7089500, 8118264, 9289091 }
	maxDegree := len(knownPartitionNumbers) - 1
	if testing.Short() {
		maxDegree = 20
	}
	for i := 1; i <= maxDegree; i++ {
		var yield chan Partition = YieldAllPartitions(i)
		var k int
		for P := range yield {
			k++
			//t.Logf("%v", P)
			sum := P.Sum()
			if sum != i {
				t.Errorf("partition check failed; sum=%d expected=%d partition=%v", sum, i, P)
			}
		}
		goodk := knownPartitionNumbers[i]
		if k != goodk {
			t.Errorf("partition number mismatch; degree=%d expected=%d got=%d", i, goodk, k)
		}
		//t.Logf("degree=%d k_d=%d", i, k)
	}
}

func intSlicesEqual(a, b []int) bool {
	if len(a) == len(b) {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
	return false
}

func TestPartitionToCycleTypeOld(t *testing.T) {
	tcase := func(partition []int, cycletype []int) {
		p := Partition(partition)
		ct := p.CycleTypeOld()
		if !intSlicesEqual(ct, cycletype) {
			t.Errorf("expected=%v got=%v", cycletype, ct)
		}
	}
	tcase([]int{1,1,1},   []int{3,0,0})
	tcase([]int{1,1,1,2}, []int{3,1,0,0,0})
}

func TestPartitionToCycleType(t *testing.T) {
	tcase := func(partition []int, cycletype []int) {
		p := Partition(partition)
		ct := make(CycleType, p.Sum())
		p.CycleType(ct)
		if !intSlicesEqual(ct, cycletype) {
			t.Errorf("expected=%v got=%v", cycletype, ct)
		}
	}
	tcase([]int{1,1,1},   []int{3,0,0})
	tcase([]int{1,1,1,2}, []int{3,1,0,0,0})
}

func TestPartitionEqual(t *testing.T) {
	type tcase struct {
		a, b Partition
		expected bool
	}
	tcases := []tcase{
		tcase{Partition{}, Partition{}, true},
		tcase{Partition{}, Partition{8}, false},
		tcase{Partition{1}, Partition{1}, true},
		tcase{Partition{1}, Partition{2}, false},
		tcase{Partition{1, 1}, Partition{2}, false},
		tcase{Partition{1, 1}, Partition{1, 1}, true},
		tcase{Partition{1, 1}, Partition{1, 2}, false},
		tcase{Partition{1}, Partition{1, 2}, false},
	}
	for i, c := range tcases {
		val := c.a.Equal(c.b)
		if val != c.expected {
			t.Errorf("Partition Equal test failed: case=%d a=%v b=%v expected=%v got=%v", i, c.a, c.b, c.expected, val)
		}
		val = c.b.Equal(c.a)
		if val != c.expected {
			t.Errorf("Partition Equal test failed: case=%d a=%v b=%v expected=%v got=%v", i, c.a, c.b, c.expected, val)
		}
	}
}

func TestSortablePartitionsEqual(t *testing.T) {
	type tcase struct {
		P, Q SortablePartitions
		expected bool
	}
	tcases := []tcase{
		tcase{SortablePartitions{}, SortablePartitions{}, true},
		tcase{SortablePartitions{Partition{}}, SortablePartitions{}, false},
		tcase{SortablePartitions{Partition{1}}, SortablePartitions{}, false},
		tcase{SortablePartitions{Partition{1}}, SortablePartitions{Partition{1}}, true},
		tcase{SortablePartitions{Partition{1, 16}}, SortablePartitions{Partition{1, 16}}, true},
		tcase{SortablePartitions{Partition{1, 16, 32}}, SortablePartitions{Partition{1, 16}}, false},
		tcase{SortablePartitions{Partition{1, 16, 32}}, SortablePartitions{Partition{1, 16, 32}}, true},
	}
	for i, c := range tcases {
		val := c.P.Equal(c.Q)
		if val != c.expected {
			t.Errorf("SortablePartitions Equal test failed: case=%d P=%v Q=%v expected=%v got=%v", i, c.P, c.Q, c.expected, val)
		}
		val = c.Q.Equal(c.P)
		if val != c.expected {
			t.Errorf("SortablePartitions Equal test failed: case=%d P=%v Q=%v expected=%v got=%v", i, c.P, c.Q, c.expected, val)
		}
	}
}

func TestPartitionsLess(t *testing.T) {
	type tcase struct {
		a, b Partition
		expected bool
		swappedExpected bool
	}
	tcases := []tcase{
		tcase{Partition{}, Partition{}, false, false},
		tcase{Partition{1}, Partition{1}, false, false},
		tcase{Partition{1}, Partition{2}, true, false},
		tcase{Partition{1, 1}, Partition{2}, true, false},
		tcase{Partition{1, 1}, Partition{1, 1}, false, false},
		tcase{Partition{1, 1}, Partition{1, 2}, true, false},
		tcase{Partition{1}, Partition{1, 2}, true, false},
	}
	for i, c := range tcases {
		result := c.a.Less(c.b)
		if result != c.expected {
			t.Errorf("a.Less(b) test failed: case=%d a=%v b=%v expected=%v got=%v", i, c.a, c.b, c.expected, result)
		}
		resultSwapped := c.b.Less(c.a)
		if resultSwapped != c.swappedExpected {
			t.Errorf("b.Less(a) test failed case=%d a=%v b=%v expected=%v got=%v", i, c.a, c.b, c.swappedExpected, resultSwapped)
		}
	}
}

func TestSortablePartitionsLess(t *testing.T) {
	type tcase struct {
		a, b Partition
		expected bool
		swappedExpected bool
	}
	tcases := []tcase{
		tcase{Partition{}, Partition{}, false, false},
		tcase{Partition{1}, Partition{1}, false, false},
		tcase{Partition{1}, Partition{2}, true, false},
		tcase{Partition{1, 1}, Partition{2}, true, false},
		tcase{Partition{1, 1}, Partition{1, 1}, false, false},
		tcase{Partition{1, 1}, Partition{1, 2}, true, false},
		tcase{Partition{1}, Partition{1, 2}, true, false},
	}
	for i, c := range tcases {
		sortable := SortablePartitions{c.a, c.b}
		result := sortable.Less(0, 1)
		if result != c.expected {
			t.Errorf("Less(a, b) test case %d failed: a=%v b=%v expected=%v got=%v", i, c.a, c.b, c.expected, result)
		}
		resultSwapped := sortable.Less(1, 0)
		if resultSwapped != c.swappedExpected {
			t.Errorf("Less(b, a) test case %d failed: a=%v b=%v expected=%v got=%v", i, c.a, c.b, c.swappedExpected, resultSwapped)
		}
	}
}

func TestSortablePartitionsLessRuleAscCompatability(t *testing.T) {
	maxDegree := 60
	if testing.Short() {
		maxDegree = 50
	}
	//t.Logf("testing d=1..%d", maxDegree)
	for d := 1; d <= maxDegree; d++ {
		var numPartitions int
		var yield chan Partition = YieldAllPartitions(d)
		var prev Partition
		for p := range yield {
			numPartitions++
			if len(prev) != 0 {
				sortable := SortablePartitions{prev, p}
				increased := sortable.Less(0, 1)
				if !increased {
					t.Errorf("expected Less(0, 1)=true; sortable=%v", sortable)
				}
				decreased := sortable.Less(1, 0)
				if decreased {
					t.Errorf("expected Less(1, 0)=false; sortable=%v", sortable)
				}
				//t.Logf("sortable=%v Less(0, 1)=%v Less(1, 0)=%v", sortable, increased, decreased)
			}
			prev = p
		}
		//t.Logf("finished testing %d partitions of partitions of d=%d", numPartitions, d)
	}
}

func TestSortOnAllPartitionsNoOp(t *testing.T) {
	maxDegree := 60
	if testing.Short() {
		maxDegree = 50
	}
	//t.Logf("testing d=1..%d", maxDegree)
	for d := 1; d <= maxDegree; d++ {
		var P SortablePartitions = AllPartitions(d)
		var Q SortablePartitions = AllPartitions(d)
		sort.Sort(Q)
		if !P.Equal(Q) {
			t.Errorf("sort of Q produced unexpected changes")
		}
	}
}

func TestSortablePartitionsSearch(t *testing.T) {
	maxDegree := 60
	if testing.Short() {
		maxDegree = 50
	}
	//t.Logf("testing d=1..%d", maxDegree)
	for d := 1; d <= maxDegree; d++ {
		var P SortablePartitions = AllPartitions(d)
		for i, p := range P {
			index, found := P.Search(p)
			if !found {
				t.Errorf("failed to find p=%v in P=%v", p, P)
			}
			if index != i {
				t.Errorf("unexpected index result from search; expected=%d got=%d", i, index)
			}
		}
	}
}
