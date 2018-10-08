// Copyright 2018 Adam Marks

package den

import (
	"testing"
)

func TestWorkerProcessPartition(t *testing.T) {
	d := 10
	exp := NewExpanderV3(d)
	exp.ensureSortedPartitions()
	exp.ensureMarkTable()
	worker := exp.newWorker(0, exp.degree)
	type tcase struct {
		p Partition
		expectedMarks int
	}
	tcases := []tcase{
		tcase{Partition{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 0},
		tcase{Partition{1, 1, 1, 1, 1, 1, 1, 1, 2}, 1},
		tcase{Partition{1, 1, 1, 1, 1, 1, 1, 3}, 1},
		tcase{Partition{1, 1, 1, 1, 1, 1, 2, 2}, 1},
		tcase{Partition{1, 1, 1, 1, 1, 1, 4}, 2},
		tcase{Partition{1, 1, 1, 1, 1, 2, 3}, 3},
		tcase{Partition{1, 1, 1, 1, 1, 5}, 1},
		tcase{Partition{1, 1, 1, 1, 2, 2, 2}, 1},
		tcase{Partition{1, 1, 1, 1, 2, 4}, 2},
		tcase{Partition{1, 1, 1, 1, 3, 3}, 1},
		tcase{Partition{1, 1, 1, 1, 6}, 3},
		tcase{Partition{1, 1, 1, 2, 2, 3}, 3},
		tcase{Partition{1, 1, 1, 2, 5}, 3},
		tcase{Partition{1, 1, 1, 3, 4}, 5},
		tcase{Partition{1, 1, 1, 7}, 1},
		tcase{Partition{1, 1, 2, 2, 2, 2}, 1},
		tcase{Partition{1, 1, 2, 2, 4}, 2},
		tcase{Partition{1, 1, 2, 3, 3}, 3},
		tcase{Partition{1, 1, 2, 6}, 3},
		tcase{Partition{1, 1, 3, 5}, 3},
		tcase{Partition{1, 1, 4, 4}, 2},
		tcase{Partition{1, 1, 8}, 3},
		tcase{Partition{1, 2, 2, 2, 3}, 3},
		tcase{Partition{1, 2, 2, 5}, 3},
		tcase{Partition{1, 2, 3, 4}, 5},
		tcase{Partition{1, 2, 7}, 3},
		tcase{Partition{1, 3, 3, 3}, 1},
		tcase{Partition{1, 3, 6}, 3},
		tcase{Partition{1, 4, 5}, 5},
		tcase{Partition{1, 9}, 2},
		tcase{Partition{2, 2, 2, 2, 2}, 1},
		tcase{Partition{2, 2, 2, 4}, 2},
		tcase{Partition{2, 2, 3, 3}, 3},
		tcase{Partition{2, 2, 6}, 3},
		tcase{Partition{2, 3, 5}, 7},
		tcase{Partition{2, 4, 4}, 2},
		tcase{Partition{2, 8}, 3},
		tcase{Partition{3, 3, 4}, 5},
		tcase{Partition{3, 7}, 3},
		tcase{Partition{4, 6}, 5},
		tcase{Partition{5, 5}, 1},
		tcase{Partition{10}, 3},
	}
	for i, c := range tcases {
		exp.resetMarkTable()
		worker.processPartition(i, c.p)
		numMarks := exp.markTable.numMarks()
		if numMarks != c.expectedMarks {
			t.Errorf("numMarks case=%d p=%v expected=%d got=%d", i, c.p, c.expectedMarks, numMarks)
		}
	}
	
}
