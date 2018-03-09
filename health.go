package health

import (
	"encoding/json"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// Status represents application health
type Status struct {
	// Application name
	Service string `json:"service"`
	// Uptime of the application
	Uptime string `json:"up_time"`
	// StartTime of the application
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

// Checker checks if the health of a service (database, external service)
type Checker interface {
	// Check the dependency status
	Check() error
}

// CheckerResult represents the health of the dependencies
type CheckerResult struct {
	// Status (UP/DOWN/TIMEOUT)
	Status string `json:"status"`
	// Error
	Error error `json:"error"`
	// ResponseTime
	ResponseTime time.Duration `json:"response_time"`
}

// Health interface
type Health interface {
	// GetStatus return the current status of the application
	GetStatus() *Status
	// RegisterChecker register a service to check their health
	RegisterChecker(name string, check Checker)
	// ServeHTTP handler for http applications
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	// Shutdown set isShutdown flag meaning the service is shutting down
	Shutdown()
}

// Options of Health instance
type Options struct {
	checkersTimeout time.Duration
}

// New returns a new Health
func New(name string, opt Options) Health {
	if opt.checkersTimeout == 0 {
		opt.checkersTimeout = time.Second
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

// GetStatus method returns the current application health status
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

func (h *health) checkersAsync() map[string]CheckerResult {
	nCheckers := len(h.checkers)
	results := map[string]CheckerResult{}

	if nCheckers == 0 {
		return results
	}

	type result struct {
		name string
		err  error
		time time.Time
	}

	ch := make(chan result, len(h.checkers))

	for n, c := range h.checkers {
		go func(name string, c Checker) {
			ch <- result{name: name, err: c.Check(), time: time.Now()}
		}(n, c)
	}

	for {
		select {
		case r := <-ch:
			resTime := time.Since(r.time)
			status := "UP"

			if r.err != nil {
				status = "DOWN"
			}

			results[r.name] = CheckerResult{Status: status, Error: r.err, ResponseTime: resTime}

			nCheckers = nCheckers - 1

			if nCheckers == 0 {
				close(ch)
				return results
			}
		case <-time.After(h.options.checkersTimeout):
			for k := range h.checkers {
				if _, ok := results[k]; !ok {
					results[k] = CheckerResult{
						Status:       "TIMEOUT",
						ResponseTime: h.options.checkersTimeout,
					}
				}
			}

			close(ch)
			return results
		}
	}
}
