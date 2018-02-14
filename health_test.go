package health

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type checker struct {
	err error
}

func (c *checker) Check() error {
	return c.err
}

func TestGetStatus(t *testing.T) {
	h := New("service_test")
	s := h.GetStatus()

	if s.Service != "service_test" {
		t.Errorf("Expect service_test and got %s", s.Service)
	}
}

func TestCheckers(t *testing.T) {
	h := New("service_test")
	h.RegisterChecker("checker1", &checker{})

	s := h.GetStatus()

	if s.HealthCheckers["checker1"] != "OK" {
		t.Errorf("Expect OK and got %s", s.HealthCheckers["checker1"])
	}
}

func TestCheckersError(t *testing.T) {
	h := New("service_test")
	h.RegisterChecker("checker1", &checker{err: errors.New("Service unreachable")})

	s := h.GetStatus()

	if s.HealthCheckers["checker1"] != "Service unreachable" {
		t.Errorf("Expect Service unreachable and got %s", s.HealthCheckers["checker1"])
	}
}

func TestServeHTTP(t *testing.T) {
	h := New("service_test")

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

func TestShutdown(t *testing.T) {
	h := New("service_test")

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
