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

// DiffMemoryStatus contains memory statistics
type DiffMemoryStatus struct {
	// TotalAlloc is cumulative bytes allocated for heap objects.
	TotalAlloc int64 `json:"total_alloc"`
	// HeapAlloc is bytes of allocated heap objects.
	HeapAlloc int64 `json:"heap_alloc"`
	// ResidentSetSize is bytes of heap memory obtained from the OS.
	ResidentSetSize int64 `json:"rss"`
}

// Memory contains the current, initial and difference statistics
type Memory struct {
	// Current statistics
	Current MemoryStatus `json:"current"`
	// Inital statistics when Health was created
	Initial MemoryStatus `json:"initial"`
	// Diff statistics between Current - Initial
	Diff DiffMemoryStatus `json:"diff"`
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

func diffMemoryStatus(current MemoryStatus, initial MemoryStatus) DiffMemoryStatus {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return DiffMemoryStatus{
		HeapAlloc:       int64(current.HeapAlloc - initial.HeapAlloc),
		TotalAlloc:      int64(current.TotalAlloc - initial.TotalAlloc),
		ResidentSetSize: int64(current.ResidentSetSize - initial.ResidentSetSize),
	}
}
