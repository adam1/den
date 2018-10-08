// Copyright 2018 Adam Marks

package main

import (
	"den"
	"flag"
	"fmt"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dpdf"
	"image/color"
	"log"
	"math/big"
	"os"
)

var bigZero = big.NewInt(0)

func main() {
	var format string
	var maxDegree int
	var outFile string

	flag.StringVar(&format, "format", "html", "html or pdf")
	flag.IntVar(&maxDegree, "max-degree", 9, "max degree of symmetric group")
	flag.StringVar(&outFile, "out", "abel-table.out", "output file")
	flag.Parse()

	tab := den.NewAbelTable(maxDegree)
	tab.Generate()

	switch (format) {
	case "html":
		generateHtml(tab, outFile)
	case "pdf":
		generatePdf(tab, outFile)
	default:
		panic(fmt.Sprintf("Unknown format type: %s", format))
	}
}

func generateHtml(tab *den.AbelTable, outFile string) {
	drawer := NewHtmlDrawer(tab, outFile)
	drawer.Draw()
	log.Printf("wrote %s", outFile)
}

type HtmlDrawer struct {
	tab *den.AbelTable
	outFile string
}

func NewHtmlDrawer(tab *den.AbelTable, outFile string) *HtmlDrawer {
	return &HtmlDrawer{
		tab: tab,
		outFile: outFile,
	}
}

// xxx todo:
// * figure out empty cols/rows
// * add checks (include non-maximal types?
// * move partition type string to row header: e.g. (3,2^3,1*)

func (h *HtmlDrawer) Draw() {
	f, err := os.Create(h.outFile)
	if err != nil {
		panic(fmt.Errorf("failed to open file for writing: %s %v", h.outFile, err))
	}
	title := fmt.Sprintf("Type Widths of Symmetric Groups up to Degree %d", h.tab.MaxDegree)

	_, err = f.WriteString(fmt.Sprintf(`
<html>
<head>
<title>%s</title>
<style TYPE="text/css">
    <!--
    h1 { font-size: 12px }
    table { font-size: 12px; border-collapse: collapse }
    table, th, td { border: 1px solid black; padding: 0.17em }
    -->
</style>
</head>
<body>
<h1>%s</h1>
<table>
`, title, title))
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(fmt.Sprintf("<tr>\n"))
	if err != nil {
		panic(err)
	}
	for i := 0; i <= len(h.tab.StringTable[0]); i++ { // note extra column for row header
		s := fmt.Sprintf("<td>%d</td>\n", i)
		if i == 0 {
			s = "<td></td>\n"
		}
		_, err = f.WriteString(s)
		if err != nil {
			panic(err)
		}
	}
	_, err = f.WriteString(fmt.Sprintf("</tr>\n"))
	if err != nil {
		panic(err)
	}
	for _, row := range h.tab.StringTable {
		_, err = f.WriteString(fmt.Sprintf("<tr>\n<td>%s</td>\n", rowHeader(row))) /// xxx first non empty string
		if err != nil {
			panic(err)
		}
		for _, s := range row {
			_, err = f.WriteString("<td>" + s + "</td>\n")
			if err != nil {
				panic(err)
			}
		}
		_, err = f.WriteString("</tr>\n")
		if err != nil {
			panic(err)
		}
	}
	
	_, err = f.WriteString(`
</table>
</body>
</html>
`)
	if err != nil {
		panic(err)
	}
	if err = f.Close(); err != nil {
		panic(err)
	}
}

func rowHeader(row []string) string {
	s := firstNonemptyItem(row)

	// xxx extract partition type label and 1-starify,
	// e.g. (6,2^2,1^2) -> (6,2^2,1*)

	// xxx maybe store in structure in node

	return s
}

func firstNonemptyItem(row []string) string {
	s := ""
	for _, x := range row {
		if x != "" {
			s = x
			break
		}
	}
	return s
}

////////////////////////////////////////////////////////////

func generatePdf(tab *den.AbelTable, outFile string) {
	var err error

	dest := draw2dpdf.NewPdf("L", "mm", "letter")
	gc := draw2dpdf.NewGraphicContext(dest)

	gc.SetFontData(draw2d.FontData{Name: "times", Family: draw2d.FontFamilySerif})

	gc.SetFontSize(2)
	gc.SetFillColor(color.Black)
	gc.SetStrokeColor(color.Black)
	gc.SetLineWidth(5)

	drawer := NewPdfDrawer(gc)
	drawer.Draw(tab)

	if err = draw2dpdf.SaveToPdfFile(outFile, dest); err != nil {
		log.Printf("Failed to write file %s: %v", outFile, err)
		os.Exit(1)
	}
}


type PdfDrawer struct {
	cellWidth, cellHeight, cellSpacing float64
	gc *draw2dpdf.GraphicContext
}

func NewPdfDrawer(gc *draw2dpdf.GraphicContext) *PdfDrawer {
	return &PdfDrawer{
		cellSpacing: 1,
		gc: gc,
	}
}

func (dr *PdfDrawer) Draw(tab *den.AbelTable) {
	dr.cellWidth, dr.cellHeight = dr.maxCellSize(tab)
	//log.Printf("xxx cellWidth=%f cellHeight=%f", dr.cellWidth, dr.cellHeight)
	//dr.drawPartitionsHeaderRow(tab)
	dr.drawTable(tab)
}

func (dr *PdfDrawer) maxCellSize(tab *den.AbelTable) (width, height float64) {
	for _, row := range tab.StringTable {
		for _, s := range row {
			_, _, r, b := dr.gc.GetStringBounds(s)
			//log.Printf("xxx s=%s l=%f t=%f r=%f b=%f", s, l, t, r, b)
			if r > width {
				width = r
			}
			if b > height {
				height = b
			}
		}
	}
	return
}

// func (dr *PdfDrawer) drawPartitionsHeaderRow(tab *den.AbelTable) {
// 	var x float64 = 10
// 	var y float64 = 10
// 	for i, p := range tab.Partitions {
// 		ct := p.CycleTypeOld()
// 		s := ct.String()
// 		dr.gc.FillStringAt(s, x, y)
// 		x += dr.cellWidth + dr.cellSpacing
// 	}
// }
func (dr *PdfDrawer) drawTable(tab *den.AbelTable) {
	var xBegin float64 = 10
	var yBegin float64 = 10
	var x float64 = xBegin
	var y float64 = yBegin
	for _, row := range tab.StringTable {
		for _, s := range row {
			dr.gc.FillStringAt(s, x, y)
			x += dr.cellWidth + dr.cellSpacing
		}
		y += dr.cellHeight
		x = xBegin
	}
}
