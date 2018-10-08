package=den
goargs=-v

all: godeps build test

default: all

run: build check-den

build:
	go build $(goargs) $(package)
	go install $(goargs) $(package)
	go install $(goargs) $(package)/check-pre-extensions
	go install $(goargs) $(package)/den
	go install $(goargs) $(package)/expander
	go install $(goargs) $(package)/gen-cpt
	go install $(goargs) $(package)/gen-partitions
	go install $(goargs) $(package)/gen-pft
	go install $(goargs) $(package)/maximal-types-matrix
	go install $(goargs) $(package)/sequence
	go install $(goargs) $(package)/abel-table

test:
	go test -short $(goargs) $(package)

godeps:
	gpm install

demo-maximal-types-matrix:
	bin/maximal-types-matrix -n 7

maximal-types.txt:
	for z in {1..20}; do bin/maximal-types-matrix -n $$z; done > $@

seq-num-maximal-types.txt:
	bin/sequence -n num_maximal_types -e 20 > $@

seq-min-cardinality-of-centralizer-of-maximal-types.txt:
	bin/sequence -n min_cardinality_centralizer_maximal_type -e 30 > $@

abel:
	go install $(goargs) $(package)/abel-table
	bin/abel-table && open hello.pdf

abel-range:
	for x in $$(seq 1 10); do bin/abel-table -max-degree $$x -out abel.$$x.html; done

clean:
	rm -rf bin pkg

.PHONY: den.txt w.txt maximal-types.txt godeps seq-min-cardinality-of-centralizer-of-maximal-types.txt profile/mem.40.pdf profile/cpu.40.pdf
