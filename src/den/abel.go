// Copyright 2018 Adam Marks

package den

import (
	"fmt"
	"log"
)

type AbelTable struct {
	MaxDegree int
	Partitions []Partition
	StringTable [][]string
}

func NewAbelTable(maxDegree int) *AbelTable {
	return &AbelTable{
		MaxDegree: maxDegree,
	}
}

// let N be the max n in the table.  we will generate all partitions
// of N; these will be displayed in the header row.  the subsequent
// rows will count from n=1 to n=N and will refer to partitions of n,
// with the specification that if row n, column j refers to a
// partition
//
//   (l_1, l_2, ..., l_m),
//
// then row n+1, column j refers to the partition
//
//   (1, l_1, l_2, ..., l_m)
//
// i.e. the same partition with the addition of a 1-cycle.
//
// of course for a given row n, there will be partitions that do not
// contain any 1-cycles, and these will be placed to the right.
//
// produce a table of strings.
//
// finally, transpose the table for better printing fit.
func (tab *AbelTable) Generate() {
	tab.Partitions = AllPartitions(tab.MaxDegree)
	tab.StringTable = make([][]string, tab.MaxDegree)
	for n := 1 ; n <= tab.MaxDegree; n++ {
		tab.generateRow(n)
	}
	log.Printf("xxx table before: %v", tab.StringTable)
	tab.StringTable = transposeStringTable(tab.StringTable)
	log.Printf("xxx table after: %v", tab.StringTable)
}

func (tab *AbelTable) generateRow(n int) {
	exp := NewExpanderV3(n)
	exp.Expand()
	partitions := AllPartitions(n)
	tab.StringTable[n - 1] = make([]string, len(tab.Partitions)) // make all rows equal to the length of the longest row
	var t CycleType = make([]int, n)
	for i, p := range partitions {
		// xxx to start, just put the type and the width in
		// the cell as a string, verify that partitions line
		// up vertically as extensions by visual inspection.
		p.CycleType(t)
		z := exp.TypeWidth(i, p, t)
		var s string
// 		if z.Cmp(bigZero) > 0 {
			s = fmt.Sprintf("%v %v", t.String(), z)
// 		}
		tab.StringTable[n - 1][i] = s
	}
}

// assumes all rows are equal length
func transposeStringTable(table [][]string) [][]string {
	width := len(table[0])
	height := 0
	for _, _ = range table {
		height += 1
	}
	transposeWidth := height
	transposeHeight := width
	transpose := make([][]string, transposeHeight)
	for i := 0; i < transposeHeight; i++ {
		transpose[i] = make([]string, transposeWidth)
		for j := 0; j < transposeWidth; j++ {
			transpose[i][j] = table[j][i]
		}
	}
	return transpose
}

