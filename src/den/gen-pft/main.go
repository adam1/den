// Copyright 2018 Adam Marks

package main

import (
	"den"
	"flag"
	"fmt"
)

func main() {

	var degree int
	flag.IntVar(&degree, "degree", 15, "degree of symmetric group")

	var lambda den.CycleType = den.CycleType{0,2,0,0,1,1}
	flag.Var(&lambda, "lambda", "cycle type in Sagan notation; comma-separated list of cycle length occurrences")

	flag.Parse()

	var P *den.PFT = den.NewPFT(degree, lambda)

	P.Generate()

	P.Check()

	fmt.Printf("P=\n%v\n", P)
}
