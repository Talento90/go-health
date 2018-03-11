package health

import "testing"

func TestNewDiffMemoryStatus(t *testing.T) {
	memStatus := MemoryStatus{
		HeapAlloc:       3,
		TotalAlloc:      5,
		ResidentSetSize: 6,
	}

	initMemStatus := MemoryStatus{
		HeapAlloc:       1,
		TotalAlloc:      2,
		ResidentSetSize: 3,
	}

	diff := diffMemoryStatus(memStatus, initMemStatus)

	if diff.TotalAlloc != 3 {
		t.Errorf("Expect TotalAlloc to be 3 and got %d", diff.TotalAlloc)
	}

	if diff.HeapAlloc != 2 {
		t.Errorf("Expect HeapAlloc to be 2 and got %d", diff.HeapAlloc)
	}

	if diff.ResidentSetSize != 3 {
		t.Errorf("Expect ResidentSetSize to be 3 and got %d", diff.ResidentSetSize)
	}
}
