package bot

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCMD(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/stock=APPL.US", "APPL.US"},
		{"/stock=GOOG.US", "GOOG.US"},
		{"/stock=MSFT.US", "MSFT.US"},
	}

	for _, test := range tests {
		result := parseCMD(test.input)
		if result != test.expected {
			t.Errorf("expected %v, got %v", test.expected, result)
		}
	}
}

func TestProcessCMD(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"help", "/help", helpMenu},
		{"appl", "/stock=AAPL.US", "AAPL.US value is $"},
		{"invalid stock", "/stock=INVALID", "INVALID stock not found; do you want to try another one?"},
		{"invalid command", "/unknown", "invalid command;\n" + helpMenu},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, _ := ProcessCMD(test.input)
			assert.True(t, strings.Contains(result.Content, test.expected))
		})
	}
}
