// Copyright 2018 Adam Marks

package main

import (
	"den"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"log"
	"math/big"
	"os"
	"time"
)

var verbose bool

func exit_on_error(f func () error) {

	var err error = f()

	if err != nil {
		os.Exit(1)
	}
}

func main() {
	var start int
	var end int
	var prof string

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&prof, "prof", "", "enabling profiling: cpu or mem")
	flag.IntVar(&start, "start", 7, "start degree of symmetric group")
	flag.IntVar(&end, "end", 7, "end degree of symmetric group")
	flag.Parse()

	switch (prof) {
	case "cpu":
		defer profile.Start().Stop()
	case "mem":
		defer profile.Start(profile.MemProfile).Stop()
	case "":
	default:
		panic(fmt.Sprintf("Uknown profile type: %s", prof))
	}

	for i := start; i <= end; i++ {
		process(i)
	}
}

func process(n int) {

	t0 := time.Now()

	var K *den.CPT = den.New_CPT(n)

	exit_on_error(func() error { return K.Generate()})

	exit_on_error(func() error { return K.Check()})

	var width *big.Int = K.Width()

	var order *big.Int = K.Order()

	var density *big.Rat = K.Density()

	totalTime := time.Since(t0)

	if verbose {
		log.Printf("%v", K)
	}

	log.Printf("n=%v width=%v order=%v density=%v types=%v ptime=%v gtime=%v wtime=%v ttime=%v\n", 
		n, width, order, density.FloatString(20), 
		K.NumCycleTypes(), int(K.PartitionTime.Seconds()), int(K.GenTime.Seconds()), int(K.WidthTime.Seconds()), int(totalTime.Seconds()))
}
