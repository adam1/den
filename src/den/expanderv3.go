// Copyright 2018 Adam Marks

package den

import (
	"fmt"
	"log"
	"math/big"
	"runtime"
	"sync"
	"time"
)

var bigZero *big.Int = big.NewInt(0)
var bigOne *big.Int = big.NewInt(1)

type ExpanderV3 struct
{
	degree int
	expanded bool
	markTable
	sortedPartitions SortablePartitions
	wg sync.WaitGroup
	width *big.Int
	workers []*expanderV3Worker

	TimeToGeneratePartitions time.Duration
	TimeToDistributeWork time.Duration
	TimeToWaitForWorkers time.Duration
	TimeToExpand time.Duration
	TimeToSumWidth time.Duration
	TimeTotalToComputeWidth time.Duration
}

type markTable []typeMark // indices correspond between markTable and sortedPartitions

type typeMark struct {
	marked bool // if marked, indicates that some other type
		    // raised to some power equals this type.

	height *big.Int  // the number of powers from 1 to the order
	                 // of the type for which the type raised to
	                 // the power equals the type itself.  equal
	                 // to the totient of the order.
}

func (m markTable) marked(index int) bool {
	return m[index].marked
}

func (m markTable) mark(index int) {
	m[index].marked = true
}

func (m markTable) reset() {
	for i := 0; i < len(m); i++ {
		m[i].marked = false
		m[i].height = nil
	}
}

func (m markTable) numMarks() (marks int) {
	for _, b := range m {
		if b.marked {
			marks++
		}
	}
	return marks
}

func (m markTable) numUnmarked() (unmarked int) {
	for _, b := range m {
		if !b.marked {
			unmarked++
		}
	}
	return unmarked
}

// note: takes ownership of height ptr
func (m markTable) setHeight(index int, height *big.Int) {
	m[index].height = height
}

func (m markTable) height(index int) *big.Int {
	return m[index].height
}

func NewExpanderV3(degree int) *ExpanderV3 {
	return &ExpanderV3{degree: degree}
}

func (exp *ExpanderV3) Degree() int {
	return exp.degree
}

func (exp *ExpanderV3) NumMaximalTypes() int {
	exp.Expand()
	return exp.markTable.numUnmarked()
}

func (exp *ExpanderV3) String() (s string) {
	return exp.dumpPartitionsAndMarks()
}

func (exp *ExpanderV3) Width() *big.Int {
	if exp.width != nil {
		return exp.width
	}
	exp.calculateWidth()
	return exp.width
}

func (exp *ExpanderV3) Density() *big.Rat {
	d := big.NewRat(1, 1)
	d.SetFrac(exp.Width(), exp.Order())
	return d
}

func (exp *ExpanderV3) Order() *big.Int {
	return Factorial(exp.degree)
}

func (exp *ExpanderV3) Expand() {
	if exp.expanded {
		return
	}
	exp.ensureSortedPartitions()
	exp.ensureMarkTable()
	log.Printf("begin expansion")
	t0 := time.Now()
	exp.spawnWorkers()
	exp.distributeWork()
	exp.closeWorkers()
	exp.TimeToExpand = time.Since(t0)
	log.Printf("expansion complete; n=%d exptime=%v", exp.degree, int(exp.TimeToExpand.Seconds()))
	exp.expanded = true
}

func (exp *ExpanderV3) calculateWidth() {
	t0 := time.Now()
	exp.Expand()
	// future: could parallelize this by fanning out partitions to workers again
	log.Printf("calculating width; n=%d", exp.degree)
	t1 := time.Now()
	width := big.NewInt(0)
	var t CycleType = make([]int, exp.degree)
	for i, p := range exp.sortedPartitions {
		if exp.marked(i) {
			continue
		}
		p.CycleType(t) // xxx optimization? avoid conversion?
		z := exp.TypeWidth(i, p, t)
		width.Add(width, z)
		//log.Printf("t=%v z=%v", &t, z)
	}
	exp.width = width
	exp.TimeToSumWidth = time.Since(t1)
	exp.TimeTotalToComputeWidth = time.Since(t0)
	log.Printf("done calculating width; n=%d width=%v swtime=%v wtime=%v", exp.degree, width,
		int(exp.TimeToSumWidth.Seconds()), int(exp.TimeTotalToComputeWidth.Seconds()))
}

func (exp *ExpanderV3) TypeWidth(i int, p Partition, t CycleType) *big.Int {
	if exp.marked(i) {
		return bigZero
	}
	z := t.CardinalityOfConjugacyClass()
	//log.Printf("xxx i=%d p=%v z=%v h=%v", i, p, z, exp.height(i))
	z.Div(z, exp.height(i))
	return z
}

func (exp *ExpanderV3) ensureSortedPartitions() {
	if len(exp.sortedPartitions) == 0 {
		exp.generateSortedPartitions()
	}
}

func (exp *ExpanderV3) generateSortedPartitions() {
	log.Printf("generating partitions; n=%d", exp.degree)
	t0 := time.Now()
	exp.sortedPartitions = AllPartitions(exp.degree)
	exp.TimeToGeneratePartitions = time.Since(t0)
	log.Printf("done generating partitions; n=%d parts=%d parttime=%v",
		exp.degree,
		len(exp.sortedPartitions),
		int(exp.TimeToGeneratePartitions.Seconds()))
}

func (exp *ExpanderV3) ensureMarkTable() {
	if len(exp.markTable) == 0 {
		exp.markTable = make([]typeMark, len(exp.sortedPartitions))
	}
}

func (exp *ExpanderV3) resetMarkTable() {
	exp.markTable.reset()
}

func (exp *ExpanderV3) spawnWorkers() {
	for i := 0; i < runtime.NumCPU(); i++ {
		exp.wg.Add(1)
		worker := exp.newWorker(i, exp.degree)
		exp.workers = append(exp.workers, worker)
	}
	log.Printf("spawned %d workers", len(exp.workers))
}

func (exp *ExpanderV3) distributeWork() {
	var i int
	var p Partition
	var v = len(exp.workers)
	t0 := time.Now()
	log.Printf("distributing to workers; workers=%d", v)
	for i, p = range exp.sortedPartitions {
		//log.Printf("i=%d p=%v work", i, p)
		exp.workers[i % v].queue <- expanderV3WorkerWork{p, i}
	}
	exp.TimeToDistributeWork = time.Since(t0)
	log.Printf("done distributing to workers; disttime=%v", int(exp.TimeToDistributeWork.Seconds()))
}

func (exp *ExpanderV3) closeWorkers() {
	for _, w := range exp.workers {
		close(w.queue)
	}
	t0 := time.Now()
	log.Printf("waiting on workers")
	exp.wg.Wait()
	exp.TimeToWaitForWorkers = time.Since(t0)
	log.Printf("done waiting for workers; wtime=%v", int(exp.TimeToWaitForWorkers.Seconds()))
	for i, w := range exp.workers {
		log.Printf("worker %d result: transactions=%d tps=%d", i,
			w.result.transactionCount, int(w.result.transactionRate))
	}
}

func (exp *ExpanderV3) dumpPartitionsAndMarks() (s string) {
	s = "types:\n"
	for i, p := range exp.sortedPartitions {
		if exp.markTable.marked(i) {
			s += "/"
		}
		//s += fmt.Sprintf("%d: %v\n", i + 1, p.StringWithoutOneCycles())
		t := p.CycleTypeOld()
		s += fmt.Sprintf("%d: %v\n", i + 1, t.StringWithCarets())
	}
	return s
}

type expanderV3Worker struct {
	index int
	degree int
	queue chan expanderV3WorkerWork
	markTable
	sortedPartitions SortablePartitions
	wg *sync.WaitGroup
	result workerResult
	ta, tb CycleType
	q Partition
	qbuf []int
}

type expanderV3WorkerWork struct {
	p Partition
	index int // index in the sortedPartitions (ruleAsc order)
}

type workerResult struct {
	transactionCount int
	transactionRate float64
}

func (exp *ExpanderV3) newWorker(index, degree int) *expanderV3Worker {
	worker := &expanderV3Worker{
		index: index,
		degree: degree,
		queue: make(chan expanderV3WorkerWork, 100),
		markTable: exp.markTable,
		sortedPartitions: exp.sortedPartitions,
		wg: &exp.wg,
		ta: make([]int, degree),
		tb: make([]int, degree),
		qbuf: make([]int, degree)}
	go worker.main()
	return worker
}

func (worker *expanderV3Worker) main() {
	defer worker.wg.Done()
	t0 := time.Now()
	var i int
	for work := range worker.queue {
		//log.Printf("worker %d i=%d p=%v", worker.index, work.index, work.p)
		worker.processPartition(work.index, work.p)
		worker.result.transactionCount++
		i++
	}
	seconds := time.Since(t0).Seconds()
	if seconds > 0 {
		worker.result.transactionRate = float64(worker.result.transactionCount) / float64(seconds)
	}
}

func (worker *expanderV3Worker) processPartition(index int, p Partition) {
	debug := false
	s := ""
	if worker.markTable.marked(index) {
		return
	}
	if debug {
		s = fmt.Sprintf("%d", index + 1)
	}
	p.CycleType(worker.ta)
	k := 2 // xxx may need to be a big int at some point...
	height := big.NewInt(1)
	if worker.ta.IsIdentity() {
		goto done
	}
	for {
		worker.ta.Power(k, worker.tb)
		if worker.ta.Equal(&worker.tb) {
			height.Add(height, bigOne)
		} else {
			worker.tb.Partition(&worker.q, worker.qbuf)
			z := worker.partitionIndex(worker.q)
			worker.markTable.mark(z)
			if debug {
				s += fmt.Sprintf(" %d", z + 1)
			}
		}
		if worker.tb.IsIdentity() {
			break
		}
		k++
	}
done:
	worker.markTable.setHeight(index, height)
	if debug {
		fmt.Print(s + "\n")
	}
}

func (worker *expanderV3Worker) partitionIndex(p Partition) int {
	index, found := worker.sortedPartitions.Search(p)
	if !found {
		panic(fmt.Errorf("failed to find partition=%v degree=%d", p, worker.degree))
	}
	return index
}
