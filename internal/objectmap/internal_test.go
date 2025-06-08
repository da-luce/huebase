package objectmap

import (
	"reflect"
	"testing"
)

// -----------------------------------------------------------------------------
// splitPath
// -----------------------------------------------------------------------------

func TestSplitPath(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"A", []string{"A"}},
		{"A.B", []string{"A", "B"}},
		{"A.B.C", []string{"A", "B", "C"}},
		{"..", []string{"", "", ""}},
		{"A..B", []string{"A", "", "B"}},
	}

	for _, tc := range tests {
		result := splitPath(tc.input)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("splitPath(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

// -----------------------------------------------------------------------------
// joinPath
// -----------------------------------------------------------------------------

func TestJoinPath(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"A"}, "A"},
		{[]string{"A", "B"}, "A.B"},
		{[]string{"A", "B", "C"}, "A.B.C"},
		{[]string{"", "", ""}, ".."},
		{[]string{"A", "", "B"}, "A..B"},
	}

	for _, tc := range tests {
		result := joinPath(tc.input)
		if result != tc.expected {
			t.Errorf("joinPath(%v) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}
