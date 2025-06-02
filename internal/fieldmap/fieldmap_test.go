package fieldmap

import (
	"reflect"
	"testing"
)

// -----------------------------------------------------------------------------
// Minimal example
// -----------------------------------------------------------------------------

// Minimal structs for testing
type Source struct {
	A *string
	B *string
}

type Dest struct {
	Group1 struct {
		A *string
		B *string
	}
	Group2 struct {
		C *string
	}
}

// Minimal YAML mapping string for test
const minimalYAML = `
A:
  - Group1.A
  - Group2.C
B:
  - Group1.B
`

func TestHasFieldByPath(t *testing.T) {
	typSource := reflect.TypeOf(Source{})
	typDest := reflect.TypeOf(Dest{})

	tests := []struct {
		typ      reflect.Type
		path     string
		expected bool
	}{
		// Base16Scheme top-level keys
		{typSource, "A", true},
		{typSource, "B", true},
		{typSource, "Z", false},

		// AbstractScheme nested keys
		{typDest, "Group1.A", true},
		{typDest, "Group1.B", true},
		{typDest, "Group2.C", true},

		// Invalid paths
		{typDest, "Group3.A", false},
		{typDest, "Group1.C", false},
		{typDest, "Group1.A.B", false},
	}

	for _, tt := range tests {
		got := hasFieldByPath(tt.typ, tt.path)
		if got != tt.expected {
			t.Errorf("hasFieldByPath(%q) = %v; want %v", tt.path, got, tt.expected)
		}
	}
}

// FIXME: this test sucks
func TestLoadMappingFromString(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		wantKeys  []string
		wantVals  map[string][]string
	}{
		{
			name: "minimal",
			yamlInput: `
A:
  - Group1.A
  - Group1.B
B:
  - Group2.C
`,
			wantKeys: []string{"A", "B"},
			wantVals: map[string][]string{
				"A": {"Group1.A", "Group1.B"},
				"B": {"Group2.C"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValues, err := LoadMappingFromString(tt.yamlInput)
			if err != nil {
				t.Fatalf("LoadMappingFromString() error = %v", err)
			}

			// Check keys presence (order not guaranteed)
			for _, wantKey := range tt.wantKeys {
				if _, found := gotValues[wantKey]; !found {
					t.Errorf("missing key %q in gotValues", wantKey)
				}
			}

			// Check values for each key
			for wantKey, wantVals := range tt.wantVals {
				gotVals, ok := gotValues[wantKey]
				if !ok {
					t.Errorf("key %q missing in gotValues", wantKey)
					continue
				}

				if len(gotVals) != len(wantVals) {
					t.Errorf("values length for key %q = %d, want %d", wantKey, len(gotVals), len(wantVals))
					continue
				}
				for i := range wantVals {
					if gotVals[i] != wantVals[i] {
						t.Errorf("value[%d] for key %q = %q, want %q", i, wantKey, gotVals[i], wantVals[i])
					}
				}
			}
		})
	}
}

func TestValidateFieldPaths(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		wantBad []string
	}{
		{
			name: "all valid paths",
			paths: []string{
				"Group1.A",
				"Group1.B",
				"Group2.C",
			},
			wantBad: nil,
		},
		{
			name: "some invalid paths",
			paths: []string{
				"Group1.A",
				"Group3",   // invalid - not in struct
				"Group3.A", // invalid - not in struct
				"Group2.C",
			},
			wantBad: []string{
				"Group3",
				"Group3.A",
			},
		},
		{
			name: "all invalid paths",
			paths: []string{
				"Foo.Bar",
				"Baz.Qux.Quux",
			},
			wantBad: []string{
				"Foo.Bar",
				"Baz.Qux.Quux",
			},
		},
		{
			name:    "empty input",
			paths:   []string{},
			wantBad: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBad := validateFieldPaths(reflect.TypeOf(Dest{}), tt.paths)

			if len(gotBad) != len(tt.wantBad) {
				t.Fatalf("unexpected number of invalid paths: got %v, want %v", gotBad, tt.wantBad)
			}

			badMap := make(map[string]struct{})
			for _, b := range gotBad {
				badMap[b] = struct{}{}
			}
			for _, want := range tt.wantBad {
				if _, found := badMap[want]; !found {
					t.Errorf("expected invalid path %q not found in results", want)
				}
			}
		})
	}
}

func TestVerifyMappingString(t *testing.T) {
	sourceType := reflect.TypeOf(Source{})
	destType := reflect.TypeOf(Dest{})

	err := VerifyMappingString(sourceType, destType, minimalYAML)
	if err != nil {
		t.Fatalf("Mapping verification failed: %v", err)
	}
}

func TestGetFieldByPath(t *testing.T) {
	str1 := "hello"
	str2 := "world"
	src := &Source{
		A: &str1,
		B: &str2,
	}

	dst := &Dest{}
	dst.Group1.A = &str1
	dst.Group1.B = &str2
	dst.Group2.C = &str2

	tests := []struct {
		name     string
		input    interface{}
		path     string
		expected interface{}
	}{
		{"Source A", src, "A", &str1},
		{"Source B", src, "B", &str2},
		{"Dest Group1.A", dst, "Group1.A", &str1},
		{"Dest Group2.C", dst, "Group2.C", &str2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFieldByPath(reflect.ValueOf(tt.input), tt.path)
			if err != nil {
				t.Fatalf("getFieldByPath failed: %v", err)
			}
			if got.Interface() != tt.expected {
				t.Errorf("got %v, want %v", got.Interface(), tt.expected)
			}
		})
	}
}

func TestSetFieldByPath(t *testing.T) {
	str1 := "foo"
	str2 := "bar"

	dst := &Dest{}
	src := &Source{}

	tests := []struct {
		name   string
		target interface{}
		path   string
		value  interface{}
	}{
		{"Set Source.A", src, "A", &str1},
		{"Set Source.B", src, "B", &str2},
		{"Set Dest.Group1.A", dst, "Group1.A", &str1},
		{"Set Dest.Group2.C", dst, "Group2.C", &str2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := reflect.ValueOf(tt.value)
			err := setFieldByPath(reflect.ValueOf(tt.target), tt.path, val)
			if err != nil {
				t.Fatalf("setFieldByPath failed: %v", err)
			}

			// Confirm by reading back
			got, err := getFieldByPath(reflect.ValueOf(tt.target), tt.path)
			if err != nil {
				t.Fatalf("getFieldByPath failed: %v", err)
			}
			if got.Interface() != tt.value {
				t.Errorf("got %v, want %v", got.Interface(), tt.value)
			}
		})
	}
}

func TestApplyMappings(t *testing.T) {
	// Set up test data
	strA := "hello"
	strB := "world"
	source := &Source{
		A: &strA,
		B: &strB,
	}
	dest := &Dest{}

	// Load the YAML mapping string
	mapping, err := LoadMappingFromString(minimalYAML)
	if err != nil {
		t.Fatalf("Failed to load YAML mapping: %v", err)
	}

	// Apply source → destination
	err = ApplySourceToDestMapping(source, dest, mapping)
	if err != nil {
		t.Fatalf("ApplySourceToDestMapping failed: %v", err)
	}

	// Validate destination values
	if dest.Group1.A == nil || *dest.Group1.A != "hello" {
		t.Errorf("Group1.A = %v, want 'hello'", deref(dest.Group1.A))
	}
	if dest.Group2.C == nil || *dest.Group2.C != "hello" {
		t.Errorf("Group2.C = %v, want 'hello'", deref(dest.Group2.C))
	}
	if dest.Group1.B == nil || *dest.Group1.B != "world" {
		t.Errorf("Group1.B = %v, want 'world'", deref(dest.Group1.B))
	}

	// Reset source and reverse the mapping
	source = &Source{}
	err = ApplyDestToSourceMapping(dest, source, mapping)
	if err != nil {
		t.Fatalf("ApplyDestToSourceMapping failed: %v", err)
	}

	// Validate restored source values
	if source.A == nil || *source.A != "hello" {
		t.Errorf("source.A = %v, want 'hello'", deref(source.A))
	}
	if source.B == nil || *source.B != "world" {
		t.Errorf("source.B = %v, want 'world'", deref(source.B))
	}
}

// Helper to safely deref strings in test output
func deref(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
