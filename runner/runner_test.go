package runner

import (
	"bytes"
	"io"
	"testing"
	"github.com/atendi9/capivara/langs"
)

// MockExec implements the Exec interface for testing
type MockExec struct {
	stdout     io.ReadCloser
	startError error
	waitError  error
}

func (m *MockExec) StdoutPipe() (io.ReadCloser, error) { return m.stdout, nil }
func (m *MockExec) Start() error                      { return m.startError }
func (m *MockExec) Wait() error                       { return m.waitError }

// Helper to create a MockExec with specific output
func newMockExec(output string) *MockExec {
	return &MockExec{
		stdout: io.NopCloser(bytes.NewBufferString(output)),
	}
}

func TestTranslate(t *testing.T) {
	tests := []struct {
		lang     langs.Lang
		key      string
		expected string
	}{
		{langs.EN_US, "success", "[SUCCESS]"},
		{langs.PT_BR, "success", "[SUCESSO]"},
		{langs.EN_US, "non_existent", ""}, 
	}

	for _, tt := range tests {
		result := translate(tt.lang, tt.key)
		if tt.expected != "" && result != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, result)
		}
	}
}

func TestProcessEvent_Go(t *testing.T) {
	r := New(langs.EN_US, nil)
	
	t.Run("Test Pass Event", func(t *testing.T) {
		event := TestEvent{
			Action:  "pass",
			Package: "main",
			Test:    "TestSomething",
			Elapsed: 0.42,
		}
		r.processEvent(event) 
	})

	t.Run("Coverage Detection", func(t *testing.T) {
		event := TestEvent{
			Action:  "output",
			Package: "main",
			Output:  "coverage: 80.0% of statements\n",
		}
		r.processEvent(event)
		if r.coverages["main"] != "coverage: 80.0% of statements" {
			t.Errorf("Coverage not captured correctly")
		}
	})
}

func TestExecuteNode_Parsing(t *testing.T) {
	tapOutput := `TAP version 13
ok 1 - should define the variable # time=10.5ms
not ok 2 - should fail this one # time=5.1ms
1..2`

	mockFn := func(cmd string, args ...string) Exec {
		return newMockExec(tapOutput)
	}

	r := New(langs.EN_US, mockFn)
	
	t.Run("Parse Node Pass", func(t *testing.T) {
		line := "ok 1 - My Test # time=10ms"
		r.processNodeEvent(line)
	})
}