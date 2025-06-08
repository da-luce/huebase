package structutil_test

import (
	"reflect"
	"testing"

	"github.com/da-luce/paletteport/internal/structutil"
)

// -----------------------------------------------------------------------------
// hasNestedFields
// -----------------------------------------------------------------------------

func TestHasNestedFieldSlice(t *testing.T) {

	testStruct := newTestStruct()

	tests := []struct {
		name      string
		path      []string
		wantFound bool
		wantValue any
	}{
		{
			name:      "Top level field A",
			path:      []string{"A"},
			wantFound: true,
			wantValue: "a",
		},
		{
			name:      "Nested field D.E",
			path:      []string{"D", "E"},
			wantFound: true,
			wantValue: "e",
		},
		// Important regression test!
		{
			name:      "Nested field D.G.H",
			path:      []string{"D", "G", "H"},
			wantFound: true,
			wantValue: nil,
		},
		{
			name:      "Deeply nested field D.G.I.J.I",
			path:      []string{"D", "G", "I", "J", "I"},
			wantFound: true,
			wantValue: "deep",
		},
		{
			name:      "Non-existent field D.G.X",
			path:      []string{"D", "G", "X"},
			wantFound: false,
		},
		{
			name:      "Invalid middle field D.X.J",
			path:      []string{"D", "X", "J"},
			wantFound: false,
		},
		{
			name:      "Partially correct path D.G.I.J.X",
			path:      []string{"D", "G", "I", "J", "X"},
			wantFound: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			found, _, val := structutil.HasNestedFieldSlice(reflect.ValueOf(testStruct), tc.path)
			if found != tc.wantFound {
				t.Errorf("expected found=%v, got %v", tc.wantFound, found)
			}
			if found {
				if val.Kind() == reflect.Ptr && val.IsNil() {
					// ok
				} else if val.Interface() != tc.wantValue {
					t.Errorf("expected value=%v, got %v", tc.wantValue, val.Interface())
				}
			}
		})
	}
}

// -----------------------------------------------------------------------------
// setNestedFields
// -----------------------------------------------------------------------------

func TestSetNestedField(t *testing.T) {
	tests := []struct {
		name      string
		path      []string
		newValue  any
		wantErr   bool
		verifyVal func(s Struct) bool
	}{
		{
			name:     "Set top-level field A",
			path:     []string{"A"},
			newValue: "updated A",
			wantErr:  false,
			verifyVal: func(s Struct) bool {
				return s.A == "updated A"
			},
		},
		{
			name:     "Set nested field D.E",
			path:     []string{"D", "E"},
			newValue: "updated E",
			wantErr:  false,
			verifyVal: func(s Struct) bool {
				return s.D.E == "updated E"
			},
		},
		{
			name:     "Set deeply nested field D.G.I.J.I",
			path:     []string{"D", "G", "I", "J", "I"},
			newValue: "deep value",
			wantErr:  false,
			verifyVal: func(s Struct) bool {
				return s.D.G.I.J.I == "deep value"
			},
		},
		{
			name:     "Set with invalid path",
			path:     []string{"D", "X", "Z"},
			newValue: "won't work",
			wantErr:  true,
			verifyVal: func(s Struct) bool {
				return true // doesn't matter, should error
			},
		},
		{
			name:     "Set with incompatible type",
			path:     []string{"A"},
			newValue: 123, // int, not string
			wantErr:  true,
			verifyVal: func(s Struct) bool {
				return s.A != "123" // should not change
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := Struct{}
			err := structutil.SetNestedField(reflect.ValueOf(&s), tc.path, reflect.ValueOf(tc.newValue))
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
			}
			if err == nil && !tc.verifyVal(s) {
				t.Errorf("field was not correctly set for path %v", tc.path)
			}
		})
	}
}
