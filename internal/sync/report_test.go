package sync_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/haul/internal/sync"
)

func TestPrintReport_AllSuccess(t *testing.T) {
	results := []sync.Result{
		{Host: "10.0.0.1", Success: true},
		{Host: "10.0.0.2", Success: true},
	}
	var buf bytes.Buffer
	sync.PrintReport(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "2 succeeded") {
		t.Errorf("expected '2 succeeded' in output, got: %s", out)
	}
	if strings.Contains(out, "failed") && !strings.Contains(out, "0 failed") {
		t.Errorf("unexpected failure mention in output: %s", out)
	}
}

func TestPrintReport_WithFailure(t *testing.T) {
	results := []sync.Result{
		{Host: "10.0.0.1", Success: true},
		{Host: "10.0.0.2", Success: false, Err: errors.New("connection refused")},
	}
	var buf bytes.Buffer
	sync.PrintReport(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "connection refused") {
		t.Errorf("expected error message in output, got: %s", out)
	}
	if !strings.Contains(out, "1 failed") {
		t.Errorf("expected '1 failed' in output, got: %s", out)
	}
}

func TestHasFailures_True(t *testing.T) {
	results := []sync.Result{
		{Host: "10.0.0.1", Success: false, Err: errors.New("err")},
	}
	if !sync.HasFailures(results) {
		t.Error("expected HasFailures to return true")
	}
}

func TestHasFailures_False(t *testing.T) {
	results := []sync.Result{
		{Host: "10.0.0.1", Success: true},
	}
	if sync.HasFailures(results) {
		t.Error("expected HasFailures to return false")
	}
}
