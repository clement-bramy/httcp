package main

import (
	"io"
	"strings"
	"testing"
)

func TestGetLinesChannel(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []string
	}{
		"empty": {
			input:    "",
			expected: []string{},
		},
		"one line": {
			input:    "hello",
			expected: []string{"hello"},
		},
		"two lines": {
			input:    "hello\nworld",
			expected: []string{"hello", "world"},
		},
	}

	for name, test := range tests {
		reader := io.NopCloser(strings.NewReader(test.input))

		actual := make([]string, 0)
		for line := range getLinesChannel(reader) {
			actual = append(actual, line)
		}

		if len(actual) != len(test.expected) {
			t.Errorf("expected lenght to be %v, but was %v", len(test.expected), len(actual))
		}

		for i, line := range actual {
			if line != test.expected[i] {
				t.Errorf("[%s] expected line %d to be: [%q], but was: [%q]", name, i, test.expected[i], line)
			}
		}
	}

}
