// Copyright 2018 Adam Marks

package main

import (
	"den"
	"flag"
	"fmt"
	"os"
)

func exit_on_error(f func () error) {

	var err error = f()

	if err != nil {
		os.Exit(1)
	}
}

func main() {

	var degree int

	flag.IntVar(&degree, "n", 7, "degree of symmetric group")

	flag.Parse()

	var cpt *den.CPT = den.New_CPT(degree)

	exit_on_error(func() error { return cpt.Generate()})

	exit_on_error(func() error { return cpt.Check()})

	cpt.Width()

	fmt.Print(cpt.MaximalTypesMatrixString())
	fmt.Print("\n")
}
