// Copyright 2018 Adam Marks

package main

import (
	"den"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"log"
	"math"
	"math/big"
	"os"
	"runtime"
)

func main() {
	begin := 1
	end := 10
	list := false
	var prof string

	flag.IntVar(&begin, "b", begin, "begin index")
	flag.IntVar(&end, "e", end, "end index")
	flag.BoolVar(&list, "l", list, "list available sequence names")
	flag.StringVar(&prof, "prof", "", "enabling profiling: cpu or mem")
	flag.Parse()

	if list {
		listSequences()
		return
	}

	switch (prof) {
	case "cpu":
		defer profile.Start().Stop()
	case "mem":
		defer profile.Start(profile.MemProfile).Stop()
	case "":
	default:
		panic(fmt.Sprintf("Uknown profile type: %s", prof))
	}

	seqNames := flag.Args()
	sequences := NewSequences(seqNames)
	printHeader(seqNames)

	for i := begin; i <= end; i++ {
		fmt.Printf("%d", i)
		for _, seq := range sequences {
			fmt.Printf(" %v", seq.ValueAtIndex(i))
		}
		fmt.Printf("\n")
	}
}

func listSequences() {
	for _, seq := range availableSequences {
		fmt.Printf("%s\n", seq.Name)
	}
}

func printHeader(seqNames []string) {
	fmt.Printf("#n")
	for _, seqName := range seqNames {
		fmt.Printf(" %s", seqName)
	}
	fmt.Printf("\n")
}

type Sequence interface {
	ValueAtIndex(n int) interface{}
}

// junk drawer
type SequenceContext struct {
	exp *den.Expander
	expV3 map[int]*den.ExpanderV3
	cpt *den.CPT
	needsPrevCpt bool
	prevCpt *den.CPT
	cumulativeDensitySum float64
}

func NewSequenceContext() *SequenceContext {
	return &SequenceContext{expV3: make(map[int]*den.ExpanderV3)}
}

func (ctx *SequenceContext) Cpt(n int) *den.CPT {
	if ctx.cpt != nil && ctx.cpt.Degree() == n {
		return ctx.cpt
	}
	if ctx.prevCpt != nil && ctx.prevCpt.Degree() == n {
		return ctx.prevCpt
	}
	if ctx.needsPrevCpt && ctx.cpt != nil {
		log.Printf("Saving CPT context for degree %d", ctx.cpt.Degree())
		ctx.prevCpt = ctx.cpt
	}
	ctx.cpt = den.New_CPT(n)
	if err := ctx.cpt.Generate(); err != nil {
		panic(err)
	}
	if err := ctx.cpt.Check(); err != nil {
		panic(err)
	}
	return ctx.cpt
}

func (ctx *SequenceContext) SetNeedsPrevCpt(b bool) {
	ctx.needsPrevCpt = b
}

func (ctx *SequenceContext) Expander(n int) *den.Expander {
	if ctx.exp == nil || ctx.exp.Degree() != n {
		ctx.exp = den.NewExpander(n)
	}
	return ctx.exp
}

func (ctx *SequenceContext) ExpanderV3(n int) *den.ExpanderV3 {
	if _, found := ctx.expV3[n]; !found {
		ctx.expV3[n] = den.NewExpanderV3(n)
	}
	return ctx.expV3[n]
}

func NewSequences(names []string) []Sequence {
	context := NewSequenceContext()
	sequences := make([]Sequence, len(names))
	for i, name := range names {
		sequences[i] = NewSequenceByName(name, context)
	}
	return sequences
}

type NamedSequenceConstructor struct {
	Name string
	Constructor func(*SequenceContext) Sequence
}

var availableSequences []*NamedSequenceConstructor = []*NamedSequenceConstructor{
	&NamedSequenceConstructor{"Density", NewDensitySequence},
	&NamedSequenceConstructor{"DensityV3", NewDensityV3Sequence},
	&NamedSequenceConstructor{"DensityDelta", NewDensityDeltaSequence},
	&NamedSequenceConstructor{"DensitySum", NewDensitySumSequence},
	&NamedSequenceConstructor{"MinCardinalityCentralizerMaximalType", NewMinCardinalityCentralizerMaximalTypeSequence},
	&NamedSequenceConstructor{"MinTotientLcmMaximalType", NewMinTotientLcmMaximalTypeSequence},
	&NamedSequenceConstructor{"NumMaximalTypes", NewNumMaximalTypesSequence},
	&NamedSequenceConstructor{"NumMaximalTypesV3", NewNumMaximalTypesV3Sequence},
	&NamedSequenceConstructor{"NumTypes", NewNumTypesSequence},
	&NamedSequenceConstructor{"TypeStoreSizeWithParts", NewTypeStoreSizeWithPartsSequence},
	&NamedSequenceConstructor{"TypeStoreSizeWithSlots", NewTypeStoreSizeWithSlotsSequence},
	&NamedSequenceConstructor{"TypeStoreSortTime", NewTypeStoreSortTimeSequence},
	&NamedSequenceConstructor{"Width", NewWidthSequence}, // xxx delete or rename old v2 stuff
	&NamedSequenceConstructor{"WidthV3", NewWidthV3Sequence},
	&NamedSequenceConstructor{"WidthV3Time", NewWidthV3TimeSequence},
	&NamedSequenceConstructor{"WidthV3SuccessiveRatio", NewWidthV3SuccessiveRatioSequence},
	&NamedSequenceConstructor{"WidthV3RatioToPreviousFactorial", NewWidthV3RatioToPreviousFactorialSequence},
	&NamedSequenceConstructor{"WidthV3RatioToPreviousFactorialTimesSquareRoot", NewWidthV3RatioToPreviousFactorialTimesSquareRootSequence}}

func NewSequenceByName(name string, context *SequenceContext) Sequence {
	var sequence *NamedSequenceConstructor
	for _, seq := range availableSequences {
		if name == seq.Name {
			sequence = seq
			break
		}
	}
	if sequence == nil {
		log.Printf("unknown sequence name: %s\n", name)
		os.Exit(1)
	}
	return sequence.Constructor(context)
}

////////////////////////////////////////////////////////////
type NumMaximalTypesSequence struct {
	context *SequenceContext
}

func NewNumMaximalTypesSequence(context *SequenceContext) Sequence {
	return &NumMaximalTypesSequence{context}
}

func (s *NumMaximalTypesSequence) ValueAtIndex(n int) interface{} {
	cpt := s.context.Cpt(n)
	return float64(cpt.NumMaximalTypes())
}

////////////////////////////////////////////////////////////
type NumTypesSequence struct {
	context *SequenceContext
}

func NewNumTypesSequence(context *SequenceContext) Sequence {
	return &NumTypesSequence{context}
}

func (s *NumTypesSequence) ValueAtIndex(n int) interface{} {
	var yield chan den.Partition = den.YieldAllPartitions(n)
	k := big.NewInt(0)
	one := big.NewInt(1)
	for range yield {
		k.Add(k, one)
	}
	return(k)
}

////////////////////////////////////////////////////////////
type MinCardinalityCentralizerMaximalTypeSequence struct {
	context *SequenceContext
}

func NewMinCardinalityCentralizerMaximalTypeSequence(context *SequenceContext) Sequence {
	return &MinCardinalityCentralizerMaximalTypeSequence{context}
}

func (s *MinCardinalityCentralizerMaximalTypeSequence) ValueAtIndex(n int) interface{} {
	cpt := s.context.Cpt(n)
	return float64(cpt.MinCardinalityCentralizerMaximalType().Int64())
}

////////////////////////////////////////////////////////////
type MinTotientLcmMaximalTypeSequence struct {
	context *SequenceContext
}

func NewMinTotientLcmMaximalTypeSequence(context *SequenceContext) Sequence {
	return &MinTotientLcmMaximalTypeSequence{context}
}

func (s *MinTotientLcmMaximalTypeSequence) ValueAtIndex(n int) interface{} {
	cpt := s.context.Cpt(n)
	return float64(cpt.MinTotientLcmMaximalType().Int64())
}

////////////////////////////////////////////////////////////
type DensitySumSequence struct {
	context *SequenceContext
}

func NewDensitySumSequence(context *SequenceContext) Sequence {
	return &DensitySumSequence{context}
}

func (s *DensitySumSequence) ValueAtIndex(n int) interface{} {
	cpt := s.context.Cpt(n)
	x, exact := cpt.Density().Float64()
	if !exact {
		log.Printf("warning: inexact Float64 n=%d x=%f", n, x)
	}
	s.context.cumulativeDensitySum += x
	return s.context.cumulativeDensitySum
}

////////////////////////////////////////////////////////////
type DensitySequence struct {
	context *SequenceContext
}

func NewDensitySequence(context *SequenceContext) Sequence {
	return &DensitySequence{context}
}

func (s *DensitySequence) ValueAtIndex(n int) interface{} {
	cpt := s.context.Cpt(n)
	den, exact := cpt.Density().Float64()
	if !exact {
		log.Printf("warning: inexact Float64 n=%d den=%f", n, den)
	}
	log.Printf("n=%d den=%v partitiontime=%d gentime=%d widthtime=%d",
		n, den, int(cpt.PartitionTime.Seconds()),
		int(cpt.GenTime.Seconds()),
		int(cpt.WidthTime.Seconds()))
	return den
}

////////////////////////////////////////////////////////////
type DensityDeltaSequence struct {
	context *SequenceContext
}

func NewDensityDeltaSequence(context *SequenceContext) Sequence {
	context.SetNeedsPrevCpt(true)
	return &DensityDeltaSequence{context}
}

func (s *DensityDeltaSequence) ValueAtIndex(n int) interface{} {
	cpt := s.context.Cpt(n)
	prevDensity := big.NewRat(0, 1)
	if n > 1 {
		prevDensity = s.context.Cpt(n-1).Density()
	}
	delta := big.NewRat(0, 1)
	delta.Sub(cpt.Density(), prevDensity)
	x, exact := delta.Float64()
	if !exact {
		log.Printf("warning: inexact Float64 n=%d x=%f", n, x)
	}
	return x
}

////////////////////////////////////////////////////////////
type WidthSequence struct {
	context *SequenceContext
}

func NewWidthSequence(context *SequenceContext) Sequence {
	return &WidthSequence{context}
}

func (s *WidthSequence) ValueAtIndex(n int) interface{} {
	cpt := s.context.Cpt(n)
	return cpt.Width()
}

////////////////////////////////////////////////////////////
type TypeStoreSizeWithSlotsSequence struct {
	context *SequenceContext
}

func NewTypeStoreSizeWithSlotsSequence(context *SequenceContext) Sequence {
	return &TypeStoreSizeWithSlotsSequence{context}
}

func (s *TypeStoreSizeWithSlotsSequence) ValueAtIndex(n int) interface{} {
	exp := s.context.Expander(n)
	runtime.GC()
	heapBefore := heapSize()
	types := exp.AllCycleTypes()
	runtime.GC()
	heapSize := heapSize() - heapBefore
	log.Printf("n=%d types=%d heap=%d countparttime=%d gentypetime=%d",
		n, len(types), heapSize, int(exp.TimeToCountPartitions.Seconds()),
		int(exp.TimeToGenerateCycleTypes.Seconds()))
	return heapSize
}

////////////////////////////////////////////////////////////
type TypeStoreSizeWithPartsSequence struct {
	context *SequenceContext
}

func NewTypeStoreSizeWithPartsSequence(context *SequenceContext) Sequence {
	return &TypeStoreSizeWithPartsSequence{context}
}

func (s *TypeStoreSizeWithPartsSequence) ValueAtIndex(n int) interface{} {
	exp := s.context.Expander(n)
	runtime.GC()
	heapBefore := heapSize()
	partitions := exp.AllPartitions()
	runtime.GC()
	heapSize := heapSize() - heapBefore
	log.Printf("n=%d parts=%d heap=%d countparttime=%d genparttime=%d",
		n, len(partitions), heapSize, int(exp.TimeToCountPartitions.Seconds()),
		int(exp.TimeToGeneratePartitions.Seconds()))
	return heapSize
}

func heapSize() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

////////////////////////////////////////////////////////////
type TypeStoreSortTimeSequence struct {
	context *SequenceContext
}

func NewTypeStoreSortTimeSequence(context *SequenceContext) Sequence {
	return &TypeStoreSortTimeSequence{context}
}

func (s *TypeStoreSortTimeSequence) ValueAtIndex(n int) interface{} {
	exp := s.context.Expander(n)
	runtime.GC()
	heapBefore := heapSize()
	types := exp.AllCycleTypes()
	runtime.GC()
	exp.DumpTypes(types)
	types = exp.SortCycleTypes()
	exp.DumpTypes(types)
	heapSize := heapSize() - heapBefore
	log.Printf("n=%d types=%d heap=%d tgen=%d tsort=%d", n, len(types), heapSize,
		int(exp.TimeToGenerateCycleTypes.Seconds()), int(exp.TimeToSortCycleTypes.Seconds()))
	return exp.TimeToSortCycleTypes.Seconds()
}

////////////////////////////////////////////////////////////
type NumMaximalTypesV3Sequence struct {
	context *SequenceContext
}

func NewNumMaximalTypesV3Sequence(context *SequenceContext) Sequence {
	return &NumMaximalTypesV3Sequence{context}
}

func (s *NumMaximalTypesV3Sequence) ValueAtIndex(n int) interface{} {
	exp := s.context.ExpanderV3(n)
	return float64(exp.NumMaximalTypes())
}

////////////////////////////////////////////////////////////
type WidthV3Sequence struct {
	context *SequenceContext
}

func NewWidthV3Sequence(context *SequenceContext) Sequence {
	return &WidthV3Sequence{context}
}

func (s *WidthV3Sequence) ValueAtIndex(n int) interface{} {
	exp := s.context.ExpanderV3(n)
	return exp.Width()
}

////////////////////////////////////////////////////////////
type WidthV3TimeSequence struct {
	context *SequenceContext
}

func NewWidthV3TimeSequence(context *SequenceContext) Sequence {
	return &WidthV3TimeSequence{context}
}

func (s *WidthV3TimeSequence) ValueAtIndex(n int) interface{} {
	exp := s.context.ExpanderV3(n)
	return int(exp.TimeTotalToComputeWidth.Seconds())
}

////////////////////////////////////////////////////////////
type WidthV3SuccessiveRatioSequence struct {
	context *SequenceContext
}

func NewWidthV3SuccessiveRatioSequence(context *SequenceContext) Sequence {
	return &WidthV3SuccessiveRatioSequence{context}
}

func (s *WidthV3SuccessiveRatioSequence) ValueAtIndex(n int) interface{} {
	r := big.NewRat(0, 1)
	if n > 1 {
		r.SetFrac(s.context.ExpanderV3(n).Width(),
			s.context.ExpanderV3(n - 1).Width())
	}
	x, _ := r.Float64()
	return x
}

////////////////////////////////////////////////////////////
type DensityV3Sequence struct {
	context *SequenceContext
}

func NewDensityV3Sequence(context *SequenceContext) Sequence {
	return &DensityV3Sequence{context}
}

func (s *DensityV3Sequence) ValueAtIndex(n int) interface{} {
	exp := s.context.ExpanderV3(n)
	x, _ := exp.Density().Float64()
	return x
}

////////////////////////////////////////////////////////////
type WidthV3RatioToPreviousFactorialSequence struct {
	context *SequenceContext
	widths []*big.Int
}

func NewWidthV3RatioToPreviousFactorialSequence(context *SequenceContext) Sequence {
	return &WidthV3RatioToPreviousFactorialSequence{
		context,
		make([]*big.Int, 1), // note: empty slot for the zeroth item
	}
}

func (s *WidthV3RatioToPreviousFactorialSequence) ValueAtIndex(n int) interface{} {
	sum := big.NewRat(0, 1)
	for r := 1; r <= n; r++ {
		width := big.NewInt(0)
		if r < len(s.widths) {
			width = s.widths[r]
		} else {
			exp := s.context.ExpanderV3(r)
			width = exp.Width()
			s.widths = append(s.widths, width)
		}
		prevFactorial := den.Factorial(r - 1)
		x := big.NewRat(0, 1)
		x.SetFrac(width, prevFactorial)
		sum.Add(sum, x)
		v, _ := x.Float64()
		log.Printf("xxx r=%d x=%v (float)x=%v sum=%v", r, x, v, sum)
	}
	z, _ := sum.Float64()
	return z
}

////////////////////////////////////////////////////////////
type WidthV3RatioToPreviousFactorialTimesSquareRootSequence struct {
	context *SequenceContext
	widths []*big.Int
}

func NewWidthV3RatioToPreviousFactorialTimesSquareRootSequence(context *SequenceContext) Sequence {
	return &WidthV3RatioToPreviousFactorialTimesSquareRootSequence{
		context,
		make([]*big.Int, 1), // note: empty slot for the zeroth item
	}
}

func (s *WidthV3RatioToPreviousFactorialTimesSquareRootSequence) ValueAtIndex(n int) interface{} {
	sum := big.NewFloat(0)
	for r := 1; r <= n; r++ {
		width := big.NewInt(0)
		if r < len(s.widths) {
			width = s.widths[r]
		} else {
			exp := s.context.ExpanderV3(r)
			width = exp.Width()
			s.widths = append(s.widths, width)
		}
		prevFactorial := big.NewFloat(0)
		prevFactorial.SetInt(den.Factorial(r - 1))

		x := big.NewFloat(0)
		x.SetInt(width)
		x.Quo(x, prevFactorial)

		g := big.NewFloat(math.Sqrt(float64(r)))
		x.Quo(x, g)

		sum.Add(sum, x)
		v, _ := x.Float64()
		log.Printf("xxx r=%d x=%v (float)x=%v sum=%v", r, x, v, sum)
	}
	z, _ := sum.Float64()
	return z
}

