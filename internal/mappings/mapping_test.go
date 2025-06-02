package mappings

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// hasFieldByPath checks if a struct type t has a nested field path like "SpecialColors.Background".
// Returns true if the full path exists, false otherwise.
func hasFieldByPath(t reflect.Type, path string) bool {
	parts := strings.Split(path, ".")

	curr := t
	for _, p := range parts {
		if curr.Kind() == reflect.Ptr {
			curr = curr.Elem()
		}
		if curr.Kind() != reflect.Struct {
			return false
		}

		field, ok := curr.FieldByName(p)
		if !ok {
			return false
		}
		curr = field.Type
	}
	return true
}

func validateFieldPaths(t reflect.Type, paths []string) (invalid []string) {
	for _, path := range paths {
		if !hasFieldByPath(t, path) {
			invalid = append(invalid, path)
		}
	}
	return
}

// loadYAMLKeysAndValuesFromString loads top-level keys and list-of-string values from a YAML string.
func loadYAMLKeysAndValuesFromString(yamlStr string) (keys []string, values [][]string, err error) {
	var raw map[string]interface{}
	err = yaml.Unmarshal([]byte(yamlStr), &raw)
	if err != nil {
		return nil, nil, err
	}

	for k, v := range raw {
		keys = append(keys, k)
		switch val := v.(type) {
		case []interface{}:
			var strs []string
			for _, i := range val {
				if s, ok := i.(string); ok {
					strs = append(strs, s)
				}
			}
			values = append(values, strs)
		case string:
			values = append(values, []string{val})
		default:
			values = append(values, []string{}) // unsupported type, empty slice
		}
	}

	return keys, values, nil
}

// verifyMappingString validates keys and values in YAML content given as a string.
// It checks that all keys exist in the source struct, and all values exist in the destination struct.
func verifyMappingString(sourceType, destType reflect.Type, yamlStr string) error {
	keys, values, err := loadYAMLKeysAndValuesFromString(yamlStr)
	if err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate that all keys exist in the source struct (e.g., Base16Scheme)
	invalidKeys := validateFieldPaths(sourceType, keys)
	if len(invalidKeys) > 0 {
		return fmt.Errorf("invalid keys in mapping: %v", invalidKeys)
	}

	// Flatten all values for validation against the destination struct (e.g., AbstractScheme)
	var allValues []string
	for _, valList := range values {
		allValues = append(allValues, valList...)
	}

	// Validate that all values exist in the destination struct
	invalidValues := validateFieldPaths(destType, allValues)
	if len(invalidValues) > 0 {
		return fmt.Errorf("invalid field paths in mapping values: %v", invalidValues)
	}

	return nil
}

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
func TestLoadYAMLKeysAndValuesFromString(t *testing.T) {
	tests := []struct {
		name      string
		yamlInput string
		wantKeys  []string
		wantVals  [][]string
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
			wantVals: [][]string{
				{"Group1.A", "Group1.B"},
				{"Group2.C"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeys, gotValues, err := loadYAMLKeysAndValuesFromString(tt.yamlInput)
			if err != nil {
				t.Fatalf("loadYAMLKeysAndValuesFromString() error = %v", err)
			}

			if len(gotKeys) != len(tt.wantKeys) {
				t.Fatalf("got %d keys, want %d keys", len(gotKeys), len(tt.wantKeys))
			}
			if len(gotValues) != len(tt.wantVals) {
				t.Fatalf("got %d values entries, want %d", len(gotValues), len(tt.wantVals))
			}

			for i, wantKey := range tt.wantKeys {
				if gotKeys[i] != wantKey {
					t.Errorf("key at index %d = %q, want %q", i, gotKeys[i], wantKey)
				}
				gotVals := gotValues[i]

				wantVals := tt.wantVals[i]
				if len(gotVals) != len(wantVals) {
					t.Errorf("values at key %q length = %d, want %d", wantKey, len(gotVals), len(wantVals))
					continue
				}
				for j := range wantVals {
					if gotVals[j] != wantVals[j] {
						t.Errorf("values[%d] at key %q = %q, want %q", j, wantKey, gotVals[j], wantVals[j])
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

	err := verifyMappingString(sourceType, destType, minimalYAML)
	if err != nil {
		t.Fatalf("Mapping verification failed: %v", err)
	}
}
