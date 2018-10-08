// Copyright 2018 Adam Marks

package den

import (
	"log"
	"sort"
	"time"
)

type Expander struct
{
	TimeToCountPartitions time.Duration
	TimeToGenerateCycleTypes time.Duration
	TimeToGeneratePartitions time.Duration
	TimeToSortCycleTypes time.Duration
	//TimeToExpand time.Duration
	//TimeToCalculateWidth time.Duration

	degree int
	markedCycleTypes SortableCycleTypes
	markedPartitions []MarkedPartition
}

type MarkedCycleType struct {
	CycleType CycleType
	Mark bool
}

type MarkedPartition struct {
	Partition Partition
	Mark bool
}

type SortableCycleTypes []MarkedCycleType

func (types SortableCycleTypes) Len() int {
	return len(types)
}

// xxx this is broken ... CycleType is currently using sagan-slot
// structure, not parts; see Partition type instead
func (types SortableCycleTypes) Less(i, j int) bool {
	// [1,1,1] < [1,1,2]
	// [1,1,1] < [1,1,1,1]
	t := types[i].CycleType
	u := types[j].CycleType
	for k, a := range t {
		if k >= len(u) {
			return false
		}
		if a > u[k] {
			return false
		}
	}
	return true
}

func (types SortableCycleTypes) Swap(i, j int) {
	t := types[i]
	types[i] = types[j]
	types[j] = t
}

func NewExpander(degree int) *Expander {
	return &Expander{degree: degree}
}

func (exp *Expander) Degree() int {
	return exp.degree
}

func (exp *Expander) AllCycleTypes() []MarkedCycleType {
	if exp.markedCycleTypes == nil {
		exp.markedCycleTypes = exp.generateAllCycleTypes()
	}
	return exp.markedCycleTypes
}

func (exp *Expander) generateAllCycleTypes() []MarkedCycleType {
	t0 := time.Now()
	k := exp.CountAllPartitions()
	types := make([]MarkedCycleType, k)
	var yield chan Partition = YieldAllPartitions(exp.degree)
	i := 0
	for p := range yield {
		types[i].CycleType = p.CycleTypeOld()
		i++
	}
	exp.TimeToGenerateCycleTypes = time.Since(t0)
	return types
}

func (exp *Expander) AllPartitions() []MarkedPartition {
	if exp.markedPartitions == nil {
		exp.markedPartitions = exp.generateAllPartitions()
	}
	return exp.markedPartitions
}

func (exp *Expander) generateAllPartitions() []MarkedPartition {
	t0 := time.Now()
	k := exp.CountAllPartitions()
	partitions := make([]MarkedPartition, k)
	var yield chan Partition = YieldAllPartitions(exp.degree)
	i := 0
	for p := range yield {
		partitions[i].Partition = p
		i++
	}
	exp.TimeToGeneratePartitions = time.Since(t0)
	return partitions
}

func (exp *Expander) CountAllPartitions() int {
	t0 := time.Now()
	var yield chan Partition = YieldAllPartitions(exp.degree)
	i := 0
	for range yield {
		i++
	}
	exp.TimeToCountPartitions = time.Since(t0)
	return i
}

func (exp *Expander) SortCycleTypes() []MarkedCycleType {
	t0 := time.Now()
	sort.Sort(exp.markedCycleTypes)
	exp.TimeToSortCycleTypes = time.Since(t0)
	return exp.markedCycleTypes
}

func (exp *Expander) DumpTypes(types []MarkedCycleType) {
	for i, t := range types {
		log.Printf("%d %v", i, t)
	}
}

// func (exp *Expander) Generate() error {
// 	exp.generateAllCycleTypes() // xxx use a markup struct wrapper here with a bool
// 	exp.sortCycleTypes()
// 	exp.expandCycleTypePowersAndMarkNonMaximals()
// 	return nil
// }

// func (exp *Expander) expandCycleTypePowersAndMarkNonMaximals() {
// 	t0 := time.Now()
// 	// xxx
// 	exp.TimeToExpand = time.Since(t0)
// }

// func (cpt *CPT) calculateWidth2() {
// 	t0 := time.Now()
// 	width := big.NewInt(0)
// 	for _, t := range cpt.markedCycleTypes {
// 		if t.Mark {
// 			continue
// 		}
// //		var z *big.Int = cpt.elementsWithType2(t.CycleType)
// 		// xxx use cardinalityOfCentralizerOfType
// 	}
// 	cpt.width = width
// 	cpt.TimeToCalculateWidth = time.Since(t0)
// }

// func (cpt *CPT) Width2() *big.Int {
// 	if cpt.width == nil {
// 		cpt.calculateWidth2()
// 	}
// 	return cpt.width
// }

