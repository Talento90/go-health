package health

import "testing"

func TestNewMemoryStatus(t *testing.T) {
	memStatus := newMemoryStatus()

	if memStatus.TotalAlloc >= memStatus.ResidentSetSize {
		t.Errorf("Expect ResidentSetSize to be bigger than %d", memStatus.TotalAlloc)
	}
}
