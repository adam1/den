// Copyright 2018 Adam Marks

package den

import (
	"fmt"
	"testing"
)

func TestTransposeStringTable(t *testing.T) {
	A := make([][]string, 3)
	A[0] = []string{"a", "b"}
	A[1] = []string{"c", "d"}
	A[2] = []string{"d", "e"}
	a := fmt.Sprintf("%v", A)
	ax := "[[a b] [c d] [d e]]"
	if (a != ax) {
		t.Errorf("expected=%s got=%s", ax, a)
	}
	B := transposeStringTable(A)
	b := fmt.Sprintf("%v", B)
	bx := "[[a c d] [b d e]]"
	if (b != bx) {
		t.Errorf("expected=%s got=%s", bx, b)
	}
}

