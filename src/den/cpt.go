// Copyright 2018 Adam Marks

package den

import (
	"fmt"
	"image"
	"image/color"
	//"log"
	"math/big"
	"strings"
	"time"
)

// conjugate power table (CPT)
type CPT struct {
	degree int
	cycleTypes []CycleType // index table; associates i to lambda_i
	cycleTypeMap map[string]int // string form of cycle type -> index
	result [][]int // indices
	markup [][]bool // structure mirrors result
	width *big.Int

	GenTime time.Duration       // xxx old
	PartitionTime time.Duration 
	WidthTime time.Duration
}

func New_CPT(degree int) *CPT {
	cpt := new(CPT)
	cpt.degree = degree
	cpt.cycleTypes = make([]CycleType, 0)
	cpt.cycleTypeMap = make(map[string]int)
	cpt.result = make([][]int, 0)
	return cpt
}

func (cpt *CPT) Degree() int {
	return cpt.degree
}

func (cpt *CPT) String() (s string) {
	s = "types:\n"
	for i, lambda := range cpt.cycleTypes {
		s += fmt.Sprintf("%v: %v\n", i+1, lambda.StringForPartitionWithoutOneCycles()) // +1 for sanity
	}
	s += "table:\n"
	for i, row := range cpt.result {
		for j, x := range row {
			if j > 0 {
				s += " "
			}
			if cpt.markup != nil && cpt.markup[i][j] {
				s+= "/"
			}
			s += fmt.Sprint(x + 1) // +1 for sanity
		}
		s += "\n"
	}
	return s
}

func (cpt *CPT) TableString() (s string) {
	for _, row := range cpt.result {
		s += "1" // for symmetry when text is centered
		for _, x := range row {
			s += " "
			s += fmt.Sprint(x + 1) // +1 for sanity
		}
		s += "\n"
	}
	return s
}

// xxx todo1: do not generate full PFTs, but rather use the
// CycleType.Power method.  xxx todo2: consider instead of building
// this CPT table, building the cycle power graph.
func (cpt *CPT) Generate() error {
	t0 := time.Now()
	//fmt.Printf("generating K_{S_%d}\n", cpt.degree)
	var yield chan Partition = YieldAllPartitions(cpt.degree)
	i := 0
	for x := range yield {
		lambda := x.CycleTypeOld()
		cpt.cycleTypes = append(cpt.cycleTypes, lambda)
		cpt.cycleTypeMap[lambda.HashKeyString()] = i
		i++
	}
	cpt.PartitionTime = time.Since(t0)
	t1 := time.Now()
	for _, lambda := range cpt.cycleTypes {
		var P *PFT = NewPFT(cpt.degree, lambda)
		P.Generate()
		//fmt.Printf("  P=\n%v\n", P)
		P.Check()
		// resolve to row of indices in result table
		row := make([]int, 0)
		for _, y := range P.data {
			row = append(row, cpt.cycleTypeMap[y.HashKeyString()])
		}
		cpt.result = append(cpt.result, row)
	}
	cpt.GenTime = time.Since(t1)
	return nil
}

func (cpt *CPT) Check() error {
	// xxx todo
	return nil
}

func (cpt *CPT) genMarkup() {
	if cpt.markup != nil {
		return
	}
	markup := make([][]bool, len(cpt.result))
	for i := 0; i < len(cpt.result); i++ {
		markup[i] = make([]bool, len(cpt.result[i]))
	}
	for i, row := range cpt.result {
		p := row[0]
		for j, x := range row {
			if j > 0 && x != p {
				markup[i][j] = true
				for z := 0; z < len(markup[x]); z++ {
					markup[x][z] = true
				}
			}
		}
	}
	cpt.markup = markup
// 	fmt.Printf("xxx K with markup = \n%v", cpt)
}

func (cpt *CPT) calculateWidth() {
	t0 := time.Now()
	if cpt.markup == nil {
		cpt.genMarkup()
	}
	width := big.NewInt(0)
	for i := range cpt.result {
		// xxx this is effectively calculating the number of
		// totatives of the LCM of the lengths of the cycle type cycles
		var x int = cpt.elementsPerTypeGroup(i)
		if x == 0 {
			//fmt.Printf("xxx i=%v x=%v\n", i+1, x) // +1 for sanity
			continue
		}
		var z *big.Int = cpt.elementsWithType(i)
		// sanity check
		bx := big.NewInt(int64(x))
		bz := big.NewInt(0)
		bz.Set(z)
		if bz.Mod(bz, bx).Sign() != 0 {
			panic(fmt.Sprintf("z=%v is not a multiple of x=%v", z, x))
		}
		y := big.NewInt(0)
		y.Set(z)
		y.Div(y, bx)
		width.Add(width, y)
		//fmt.Printf("xxx i=%v x=%v z=%v y=%v\n", i+1, x, z, y) // +1 for sanity
	}
	cpt.width = width
	cpt.WidthTime = time.Since(t0)
}

func (cpt *CPT) elementsPerTypeGroup(x int) int {
	if cpt.markup == nil {
		return -1
	}
	row := cpt.result[x]
	var z int
	for j := range row {
		if cpt.markup[x][j] {
			if j == 0 {
				break
			}
		} else {
			z++
		}
	}
	return z
}

func (cpt *CPT) elementsWithType(x int) *big.Int {
	var k *big.Int = Factorial(cpt.degree)
	var lambda CycleType = cpt.cycleTypes[x]
	// xxx refactor to use cardinalityOfCentralizerOfType
	for i := 0; i < cpt.degree; i++ {
		s := i + 1
		m := lambda[i]
		z := Exp(s, m)
		v := Factorial(m)
		z.Mul(z, v)
		//fmt.Printf("  xxx s=%v m=%v k=%v v=%v z=%v\n", s, m, k, v, z)
		// sanity check
		b := big.NewInt(0)
		b.Set(k)
		if b.Mod(b, z).Sign() != 0 {
			panic(fmt.Sprintf("k=%v is not a multiple of z=%v", k, z))
		}

		k.Div(k, z)
	}
	return k
}

func (cpt *CPT) CardinalityOfCentralizerOfType(t CycleType) *big.Int {
	k := big.NewInt(1)
	for i := 0; i < cpt.degree; i++ {
		s := i + 1
		m := t[i]
		k.Mul(k, Exp(s, m))
		k.Mul(k, Factorial(m))
	}
	return k
}

func (cpt *CPT) Width() *big.Int {
	if cpt.width == nil {
		cpt.calculateWidth()
	}
	return cpt.width
}

func (cpt *CPT) Density() *big.Rat {
	d := big.NewRat(1, 1)
	d.SetFrac(cpt.Width(), cpt.Order())
	return d
}

func (cpt *CPT) Order() *big.Int {
	return Factorial(cpt.degree)
}

func (cpt *CPT) NumCycleTypes() int {
	return len(cpt.cycleTypes)
}

func (cpt *CPT) Diameter() int {
	var max int
	for _, lambda := range cpt.cycleTypes {
		d := lambda.Order()
		if d > max {
			max = d
		}
	}
	return max
}

func (cpt *CPT) MaximalTypesMatrixString() string {
	var s string
	for i := 0; i < cpt.degree; i++ {
		if i > 0 {
			s += " "
		}
		s += fmt.Sprintf("%4d", i+1)
	}
	s += "\n"
	for i := 0; i < cpt.degree; i++ {
		if i > 0 {
			s += "-"
		}
		s += fmt.Sprint("----")
	}
	s += "\n"
	types := cpt.MaximalTypes()
	for _, z := range types {
		for j, t := range *z {
			if j > 0 {
				s += " "
			}
			if t > 0 {
				s += fmt.Sprintf("%4d", t)
			} else {
				s += "    "
			}
		}
		s += "  [" + z.StringForPartitionWithoutOneCycles() + "]\n"
	}
	return s
}

func (cpt *CPT) NumMaximalTypes() int {
	return len(cpt.MaximalTypes())
}

func (cpt *CPT) MaximalTypes() []*CycleType {
	cpt.Width()
	result := make([]*CycleType, 0)
	for i, z := range cpt.cycleTypes {
		if cpt.markup[i][0] { // not maximal
			continue
		}
		result = append(result, &z)
	}
	return result
}

// xxx is this correct?
func (cpt *CPT) MinCardinalityCentralizerMaximalType() *big.Int {
	m := big.NewInt(0)
	types := cpt.MaximalTypes()
	for i, z := range types {
		c := cpt.CardinalityOfCentralizerOfType(*z)
		//fmt.Printf("type [%s] c=%d\n", z.FriendlyString(), c)
		if i == 0 || c.Cmp(m) < 0 {
			m = c
		}
	}
	return m
}

func (cpt *CPT) MinTotientLcmMaximalType() *big.Int {
	// xxx todo
	return big.NewInt(0)
}

type Logarithm struct {
	Base *CycleType
	Power int
}

// Note that we only return positive powers, so we do not include all
// types to the zeroth power.
func (cpt *CPT) Logarithms(u *CycleType) []Logarithm {
	logarithms := make([]Logarithm, 0)
	key := u.HashKeyString() // todo: improve this
	uIndex := cpt.cycleTypeMap[key]
	//log.Printf("xxx Logarithm key=%v uIndex=%v", key, uIndex)
	for _, row := range cpt.result {
		baseTypeIndex := 0
		//log.Printf("xxx Logarithm check row=%v", row)
		for j, x := range row {
			//log.Printf("xxx Logarithm check j=%v x=%v", j, x)
			if j == 0 {
				baseTypeIndex = x
			}
			if x == uIndex {
				//log.Printf("xxx match uIndex=%v", uIndex)
				logarithms = append(logarithms, Logarithm{&cpt.cycleTypes[baseTypeIndex], j+1})
			}
		}
	}
	//log.Printf("xxx Logarithms returning: %v", logarithms)
	return logarithms
}

func (cpt *CPT) Latex() (s string) {
	var colstring string
	var d = cpt.Diameter()
	for i := 0; i < d; i++ {
		colstring += "r "
	}
	s = fmt.Sprintf("\\begin{tabular}{%v}\n", colstring)
	for i, row := range cpt.result {
		for j, x := range row {
			if j > 0 {
				s += " & "
			}
			z := fmt.Sprint(x + 1) // +1 for sanity
			if cpt.markup != nil && cpt.markup[i][j] {
				s += fmt.Sprintf("\\cancel{%v}", z)
			} else {
				s += z
			}
		}
		s += " \\\\\n"
	}	
	s += "\\end{tabular}\n"
	return s
}

func (cpt *CPT) Image() (img image.Image) {
	//return cpt.ImageWithCellDimensions(1, 1, 0, 0)
	return cpt.ImageWithCellDimensions(4, 4, 1, 1)
}

func (cpt *CPT) ImageWithCellDimensions(cellwidth, cellheight, paddingx, paddingy int) (img image.Image) {
	generator := &cptImageGenerator{CellWidth:cellwidth, 
		CellHeight:cellheight,
		PaddingX:paddingx,
		PaddingY:paddingy}
	return generator.generate(cpt)
}

// graphviz 
func (cpt *CPT) Dot() string {
	out := "digraph G\n"
	out += "{\n"
	table := strings.Replace(cpt.TableString(), "\n", "\\n", -1)
 	out += fmt.Sprintf("  graph [label=\"T_%d\n%v\" labelfontsize=0.5 labeldistance=10];\n", cpt.degree, table)
//	out += "node [shape=point width=0.015];\n"
//	out += "node [penwidth=0.25];\n"
	out += "node [shape=none];\n"
 	out += "edge [penwidth=0.25 arrowsize=0.5 color=\"#ff000077\"];\n"
	// vertices
	for i, t := range cpt.cycleTypes {
		out += fmt.Sprintf("\"%v\"\n", cpt.CycleTypeDescription(i, &t))
	}
	// xxx this is broken.  what is wanted is a hasse diagram.  in
	// order to make a hasse diagram, start with the full power
	// graph, then remove edges that bypass a "cover".  for now,
	// just draw the whole power graph.
	for _, row := range cpt.result {
		b := row[0]
		for j := 1; j < len(row); j++ {
			if row[j] == b {
				continue
			}
			out += fmt.Sprintf("\"%v\" -> \"%v\"\n",
				cpt.CycleTypeDescriptionFromIndex(b),
				cpt.CycleTypeDescriptionFromIndex(row[j]))
		}
	}
	out += "}\n";
	return out
}

func (cpt *CPT) CycleTypeDescriptionFromIndex(index int) string {
	return cpt.CycleTypeDescription(index, &cpt.cycleTypes[index])
}

func (cpt *CPT) CycleTypeDescription(index int, t *CycleType) string {
	return fmt.Sprintf("%d: %v", index+1, t.StringWithTilde())
}

func (cpt *CPT) LatexTypes() (s string) {
	// xxx
	return s
}

// todo: move these
type cptImageGenerator struct {
	CellWidth int
	CellHeight int
	PaddingX int
	PaddingY int
	img *image.RGBA
}

func (gen *cptImageGenerator) generate(cpt *CPT) *image.RGBA {
	cpt.genMarkup()
	diameter := cpt.Diameter()
	totalwidth := gen.CellWidth * diameter + (diameter - 1) * gen.PaddingX
	rows := len(cpt.result)
	totalheight := gen.CellHeight * rows + (rows - 1) * gen.PaddingY
	gen.img = image.NewRGBA(image.Rect(0, 0, totalwidth, totalheight))
	for i, row := range cpt.result {
		for j := 0; j < diameter; j++ {
			if j < len(row) {
				if cpt.markup[i][j] {
					gen.drawMarkedCell(i, j)
				} else {
					gen.drawUnmarkedCell(i, j)
				}
			} else {
				gen.drawBackgroundCell(i, j)
			}
		}
	}
	return gen.img
}

func (gen *cptImageGenerator) drawMarkedCell(row, col int) {
	gen.drawBackgroundCell(row, col)
}

func (gen *cptImageGenerator) drawUnmarkedCell(row, col int) {
	gen.drawCell(row, col, color.Black)
}

func (gen *cptImageGenerator) drawBackgroundCell(row, col int) {
	gen.drawCell(row, col, color.White)
}

func (gen *cptImageGenerator) drawCell(row, col int, c color.Color) {
	var xoffset int
	if col > 0 {
		xoffset = col * gen.CellWidth + (col - 1) * gen.PaddingX
	}
	var yoffset int
	if row > 0 {
		yoffset = row * gen.CellHeight + (row - 1) * gen.PaddingY
	}
	for i := 0; i < gen.CellWidth; i++ {
		for j := 0; j < gen.CellHeight; j++ {
			gen.img.Set(xoffset + i, yoffset + j, c)
		}
	}
}
