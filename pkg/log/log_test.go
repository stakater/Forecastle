package log

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewReturnsSharedLogger(t *testing.T) {
	first, second := New(), New()
	if first != second {
		t.Fatal("expected New to return the same logger on every call")
	}
}

// Logging alternately from loggers obtained in different packages must not keep
// wrapping the formatter (https://github.com/stakater/Forecastle/issues/589).
func TestAlternatingLoggersDoNotNestFormatter(t *testing.T) {
	a, b := New(), New()

	var buf bytes.Buffer
	originalOut := a.Out
	a.SetOutput(&buf)
	defer a.SetOutput(originalOut)

	a.Info("warm up")
	b.Info("warm up")
	formatterA, formatterB := a.Formatter, b.Formatter

	for i := 0; i < 50; i++ {
		a.Info("a")
		b.Info("b")
	}

	if a.Formatter != formatterA || b.Formatter != formatterB {
		t.Fatal("formatter was wrapped again after alternating between loggers")
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	last := lines[len(lines)-1]
	if !strings.Contains(last, `_source="log/log_test.go:`) {
		t.Fatalf("expected the filename hook to add a source field, got %q", last)
	}
}
