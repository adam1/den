// Copyright 2018 Adam Marks

package main

import (
	"den"
	"flag"
	"fmt"
	"os"
)

func main() {

	var degree int
	var matrixStyle bool

	flag.IntVar(&degree, "n", 6, "integer degree")

	flag.BoolVar(&matrixStyle, "m", false, "(sagan) matrix style")

	flag.Parse()

	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}
	var yield chan den.Partition = den.YieldAllPartitions(degree)

	fmt.Printf("Lambda_%d = [\n", degree)

	k := 0
	for p := range yield {
		k++
		var s string
		if matrixStyle {
			t := p.CycleTypeOld()
			s = fmt.Sprintf("%v", t.String())
		} else {
			s = p.String()
		}
		fmt.Printf("  %v\n", s)
	}
	fmt.Printf("]\nk_%d = %d", degree, k)
}
