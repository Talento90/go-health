package health

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockChecker struct {
	err       error
	sleepTime time.Duration
}

func (mc *mockChecker) Check() error {
	if mc.sleepTime > 0 {
		time.Sleep(mc.sleepTime)
	}

	return mc.err
}

func TestGetStatus(t *testing.T) {
	h := New("service_test", Options{})
	s := h.GetStatus()

	if s.Service != "service_test" {
		t.Errorf("Expect service_test and got %s", s.Service)
	}
}

func TestCheckers(t *testing.T) {
	h := New("service_test", Options{})
	h.RegisterChecker("checker1", &mockChecker{})

	s := h.GetStatus()

	if s.HealthCheckers["checker1"].Status != "CHECKED" {
		t.Errorf("Expect CHECKED and got %s", s.HealthCheckers["checker1"].Status)
	}
}

func TestCheckersError(t *testing.T) {
	h := New("service_test", Options{})
	h.RegisterChecker("checker1", &mockChecker{err: errors.New("Service unreachable")})

	s := h.GetStatus()

	if s.HealthCheckers["checker1"].Error.Error() != "Service unreachable" {
		t.Errorf("Expect Service unreachable and got %s", s.HealthCheckers["checker1"].Error)
	}
}

func TestMultipleCheckers(t *testing.T) {
	results := []string{"CHECKED", "CHECKED", "TIMEOUT", "CHECKED"}

	h := New("service_test", Options{CheckersTimeout: time.Second * 1})

	h.RegisterChecker("checker1", &mockChecker{})
	h.RegisterChecker("checker2", &mockChecker{sleepTime: time.Millisecond * 300})
	h.RegisterChecker("checker3", &mockChecker{sleepTime: time.Second * 5})
	h.RegisterChecker("checker4", &mockChecker{err: errors.New("Error connections to db"), sleepTime: time.Millisecond * 500})

	s := h.GetStatus()

	for i, expected := range results {
		cs := s.HealthCheckers[fmt.Sprintf("checker%d", i+1)].Status

		if cs != expected {
			t.Errorf("Expect checker%d to be %s and got %s", i+1, expected, cs)
		}
	}
}

func TestServeHTTP(t *testing.T) {
	h := New("service_test", Options{})

	req, err := http.NewRequest("GET", "localhost/health", nil)

	if err != nil {
		t.Error(err)
	}

	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Result().StatusCode != 200 {
		t.Errorf("Expect 200 and got %d", rec.Result().StatusCode)
	}

	s := Status{}

	bytes, err := ioutil.ReadAll(rec.Body)

	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(bytes, &s)

	if err != nil {
		t.Error(err)
	}

	if s.Service != "service_test" {
		t.Errorf("Expect service_test and got %s", s.Service)
	}
}

func TestServeHTTPNotFound(t *testing.T) {
	h := New("service_test", Options{})

	req, err := http.NewRequest("PUT", "localhost/health", nil)

	if err != nil {
		t.Error(err)
	}

	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Result().StatusCode != 404 {
		t.Errorf("Expect 404 and got %d", rec.Result().StatusCode)
	}
}

func TestShutdown(t *testing.T) {
	h := New("service_test", Options{})

	h.Shutdown()

	req, err := http.NewRequest("GET", "localhost/health", nil)

	if err != nil {
		t.Error(err)
	}

	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expect %d and got %d", http.StatusServiceUnavailable, rec.Result().StatusCode)
	}
}

func TestServeHTTPCheckerFailure(t *testing.T) {
	h := New("service_test", Options{})
	h.RegisterChecker("checker1", &mockChecker{err: errors.New("Service unreachable")})

	req, err := http.NewRequest("GET", "localhost/health", nil)

	if err != nil {
		t.Error(err)
	}

	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Result().StatusCode != 503 {
		t.Errorf("Expect 503 and got %d", rec.Result().StatusCode)
	}
}
