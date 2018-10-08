// Copyright 2018 Adam Marks

package main

import (
	"den"
	"flag"
	"fmt"
)

func main() {
	var degree int

	flag.IntVar(&degree, "n", 7, "degree of symmetric group")
	flag.Parse()

	exp := den.NewExpanderV3(degree)
	exp.NumMaximalTypes()

	fmt.Print(exp.String())
}
