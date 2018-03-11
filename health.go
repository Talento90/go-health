package health

import (
	"encoding/json"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// Status of the service
type Status struct {
	// Service name
	Service string `json:"service"`
	// Uptime of the service (How long service is running)
	Uptime string `json:"up_time"`
	// StartTime
	StartTime string `json:"start_time"`
	// Memory statistics
	Memory Memory `json:"memory"`
	// GoRoutines being used
	GoRoutines int `json:"go_routines"`
	// IsShuttingDown is active
	IsShuttingDown bool `json:"is_shutting_down"`
	// HealthCheckers status
	HealthCheckers map[string]CheckerResult `json:"health_checkers"`
}

// Checker interface checks the health of external services (database, external service)
type Checker interface {
	// Check service health
	Check() error
}

// CheckerResult struct represent the result of a checker
type CheckerResult struct {
	name string
	// Status (CHECKED/TIMEOUT)
	Status string `json:"status"`
	// Error
	Error error `json:"error"`
	// ResponseTime
	ResponseTime string `json:"response_time"`
}

// Health interface
type Health interface {
	// GetStatus return the current status of the service
	GetStatus() *Status
	// RegisterChecker register a service to check their health
	RegisterChecker(name string, check Checker)
	// ServeHTTP handler for http services
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	// Shutdown set isShutdown flag meaning the service is shutting down
	Shutdown()
}

// Options of Health instance
type Options struct {
	// CheckersTimeout is the timeout value when checking the health of registed checkers
	CheckersTimeout time.Duration
}

// New returns a new Health
func New(name string, opt Options) Health {
	if opt.CheckersTimeout == 0 {
		opt.CheckersTimeout = time.Second
	}

	return &health{
		name:       name,
		mutex:      &sync.Mutex{},
		startTime:  time.Now(),
		initMem:    newMemoryStatus(),
		checkers:   map[string]Checker{},
		isShutdown: false,
		options:    opt,
	}
}

// RegisterChecker register an external dependencies health
func (h *health) RegisterChecker(name string, check Checker) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.checkers[name] = check
}

// Shutdown the health monitor
func (h *health) Shutdown() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.isShutdown = true
}

// GetStatus method returns the current service health status
func (h *health) GetStatus() *Status {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	numGoRoutines := runtime.NumGoroutine()
	memStatus := newMemoryStatus()
	diffMemStatus := MemoryStatus{
		HeapAlloc:       memStatus.HeapAlloc - h.initMem.HeapAlloc,
		TotalAlloc:      memStatus.TotalAlloc - h.initMem.TotalAlloc,
		ResidentSetSize: memStatus.ResidentSetSize - h.initMem.ResidentSetSize,
	}

	results := h.checkersAsync()

	return &Status{
		Service:    h.name,
		Uptime:     time.Since(h.startTime).String(),
		StartTime:  h.startTime.Format(time.RFC3339),
		GoRoutines: numGoRoutines,
		Memory: Memory{
			Initial: h.initMem,
			Current: memStatus,
			Diff:    diffMemStatus,
		},
		IsShuttingDown: h.isShutdown,
		HealthCheckers: results,
	}
}

// ServeHTTP that returns the health status
func (h *health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	code := http.StatusOK
	status := h.GetStatus()
	bytes, _ := json.Marshal(status)

	if h.isShutdown {
		code = http.StatusServiceUnavailable
	}

	w.WriteHeader(code)
	w.Write(bytes)
}

type health struct {
	mutex      *sync.Mutex
	name       string
	isShutdown bool
	startTime  time.Time
	initMem    MemoryStatus
	checkers   map[string]Checker
	options    Options
}

func check(ch chan<- CheckerResult, timeout time.Duration, name string, c Checker) {
	start := time.Now()
	done := make(chan error)

	go func() {
		done <- c.Check()
	}()

	select {
	case r := <-done:
		ch <- CheckerResult{name: name, Status: "CHECKED", Error: r, ResponseTime: time.Since(start).String()}
	case <-time.After(timeout):
		ch <- CheckerResult{name: name, Status: "TIMEOUT", ResponseTime: time.Since(start).String()}
	}

}

func (h *health) checkersAsync() map[string]CheckerResult {
	numCheckers := len(h.checkers)
	results := map[string]CheckerResult{}

	if numCheckers == 0 {
		return results
	}

	ch := make(chan CheckerResult, numCheckers)

	for n, c := range h.checkers {
		go check(ch, h.options.CheckersTimeout, n, c)
	}

	i := numCheckers

	for r := range ch {
		results[r.name] = r

		i = i - 1

		if i == 0 {
			close(ch)
		}
	}

	return results
}
