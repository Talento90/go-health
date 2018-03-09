package health

import "runtime"

// MemoryStatus contains memory statistics
type MemoryStatus struct {
	// TotalAlloc is cumulative bytes allocated for heap objects.
	TotalAlloc uint64 `json:"total_alloc"`
	// HeapAlloc is bytes of allocated heap objects.
	HeapAlloc uint64 `json:"heap_alloc"`
	// ResidentSetSize is bytes of heap memory obtained from the OS.
	ResidentSetSize uint64 `json:"rss"`
}

// Memory contains the current, initial and difference statistics
type Memory struct {
	Current MemoryStatus `json:"current"`
	Initial MemoryStatus `json:"initial"`
	Diff    MemoryStatus `json:"diff"`
}

func newMemoryStatus() MemoryStatus {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return MemoryStatus{
		HeapAlloc:       mem.HeapAlloc,
		TotalAlloc:      mem.TotalAlloc,
		ResidentSetSize: mem.HeapSys,
	}
}
