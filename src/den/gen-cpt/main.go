// Copyright 2018 Adam Marks

package main

import (
	"bufio"
	"den"
	"flag"
	"fmt"
	"image/png"
	"io/ioutil"
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
	var latexCPT bool
	var latexTypes bool
	var pngOut string
	var dotOut string

	flag.IntVar(&degree, "n", 7, "degree of symmetric group")

	flag.BoolVar(&latexCPT, "latex-cpt", false, "generate latex output of cpt table")

	flag.BoolVar(&latexTypes, "latex-types", false, "generate latex output of types table")

	flag.StringVar(&pngOut, "png", "", "generate PNG file output of cpt table")

	flag.StringVar(&dotOut, "dot", "", "generate graphviz DOT file visualization of cpt table")

	flag.Parse()

	var cpt *den.CPT = den.New_CPT(degree)

	exit_on_error(func() error { return cpt.Generate()})

	exit_on_error(func() error { return cpt.Check()})

	cpt.Width()

	if pngOut != "" {
		
		img := cpt.Image()

		file, err := os.Create(pngOut)

		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)

		err = png.Encode(writer, img)

		if err != nil {
			panic(err)
		}

		err = writer.Flush()

		if err != nil {
			panic(err)
		}

	} else if dotOut != "" {

		var dot string = cpt.Dot()

		err := ioutil.WriteFile(dotOut, []byte(dot), 0644)

		if err != nil {
			panic(err)
		}

	} else if latexCPT  {

		fmt.Printf("%v", cpt.Latex())

	} else if latexTypes {

		fmt.Printf("%v", cpt.LatexTypes())

	} else {
		fmt.Printf("cpt=\n%v\n", cpt)
	}
}
