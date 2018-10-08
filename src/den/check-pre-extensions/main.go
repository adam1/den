// Copyright 2018 Adam Marks

package main

import (
	"den"
	"flag"
	"log"
	"os"
)

func main() {
	begin := 1
	end := 10

	flag.IntVar(&begin, "b", begin, "begin index")
	flag.IntVar(&end, "e", end, "end index")
	flag.Parse()

	log.Printf("Checking for non-maximal pre-extensions from n=%d to n=%d", begin, end)
	var prevCpt, cpt *den.CPT
	for i := begin; i <= end; i++ {
		log.Printf("Generating n=%d", i)
		cpt = den.New_CPT(i)
		if err := cpt.Generate(); err != nil {
			panic(err)
		}
		checkDegree(i, cpt, prevCpt)
		prevCpt = cpt
	}
}

func checkDegree(index int, cpt, prevCpt *den.CPT) {
	if prevCpt == nil {
		return
	}
	for _, t := range cpt.MaximalTypes() {
		for _, p := range t.PreExtensions() {
			log.Printf("Checking type t=%v pre-extension p=%v", t, p)
			logarithms := prevCpt.Logarithms(p)
			for _, x := range logarithms {
				if x.Power != 1 {
					log.Printf("Found logarithm for pre-extension! t=%v p=%v logs=%v", t, p, logarithms)
					os.Exit(1)
				}
			}
		}
	}
	log.Printf("Done checking n=%d maximal_types=%d", index, cpt.NumMaximalTypes())
}
