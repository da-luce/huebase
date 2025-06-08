package adapter

import (
	"fmt"
	"strings"
	"testing"

	"reflect"

	"github.com/da-luce/paletteport/internal/structutil"
)

func TestAllAdapters(t *testing.T) {
	for _, ad := range Adapters {
		t.Run(ad.Name(), func(t *testing.T) {
			t.Run("testFillDummy", func(t *testing.T) {
				testFillDummy(t, ad)
			})
			t.Run("TransitiveProperty", func(t *testing.T) {
				testTransitiveProperty(t, ad, 0.9)
			})
			t.Run("TestRenderAdapterToString", func(t *testing.T) {
				checkRenderAdapterToString(t, ad)
			})
			t.Run("TransitivePropertyPart2", func(t *testing.T) {
				testTransitivePropertyPart2(t, ad)
			})
		})
	}
}

// Helper: check all pointer fields are non-nil
func checkAllFieldsSet(t *testing.T, v interface{}) {
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			t.Errorf("field %s is nil", val.Type().Field(i).Name)
		}
	}
}

func testFillDummy(t *testing.T, ad Adapter) {
	fillDummyScheme(ad)
	checkAllFieldsSet(t, ad)
}

func testTransitiveProperty(t *testing.T, ad Adapter, simThresh float64) {
	// Create new instance of adapter type
	schemeVal := reflect.New(reflect.TypeOf(ad).Elem())
	scheme := schemeVal.Interface().(Adapter)

	// Fill dummy data so all pointer fields are non-nil
	fillDummyScheme(scheme)

	// Create another new instance of adapter type for output
	newSchemeVal := reflect.New(reflect.TypeOf(ad).Elem())
	newScheme := newSchemeVal.Interface().(Adapter)

	// Run the adapter mapping: scheme -> abstract -> newScheme
	err := adaptScheme(scheme, newScheme) // <-- no & here
	if err != nil {
		t.Fatalf("Failed to adapt scheme transitive property for %T: %v", scheme, err)
	}

	// Check that all pointer fields in newScheme are non-nil
	checkAllFieldsSet(t, newScheme)

	// Check similarity between original and new scheme
	sim := FieldSimilarity(scheme, newScheme)
	if sim < simThresh {
		t.Errorf("Similarity too low for %T: got %.2f, want >= %.2f", scheme, sim, simThresh)
	}
}

// A helper function that prints struct with pointer fields dereferenced:
// Updated printSchemeContents using TraverseStructDFS
func printSchemeContents(prefix string, scheme interface{}) {
	fmt.Printf("%s:\n", prefix)

	structutil.TraverseStructDFS(scheme, func(path []string, field reflect.StructField, value reflect.Value) bool {
		// Skip unexported fields
		if field.PkgPath != "" {
			return true
		}

		indent := strings.Repeat("\t", len(path)-1)
		fieldPath := strings.Join(path, ".")

		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				fmt.Printf("%s%s: <nil>\n", indent, fieldPath)
				return true
			}

			if isColor(value.Type()) {
				fmt.Printf("%s%s: %v\n", indent, fieldPath, value.Elem().Interface())
				return false
			}

			value = value.Elem()
		}

		if isColor(value.Type()) {
			fmt.Printf("%s%s: %v\n", indent, fieldPath, value.Interface())
			return false
		}

		fmt.Printf("%s%s: %v\n", indent, fieldPath, value.Interface())
		return true
	})
}

func checkRenderAdapterToString(t *testing.T, ad Adapter) {
	schemeVal := reflect.New(reflect.TypeOf(ad).Elem())
	scheme := schemeVal.Interface().(Adapter)

	fillDummyScheme(scheme)

	output, err := RenderAdapterToString(scheme)
	if err != nil {
		t.Fatalf("RenderAdapterToString failed for %T: %v", scheme, err)
	}
	if len(output) == 0 {
		t.Errorf("expected non-empty output string for %T", scheme)
	}
}

func testTransitivePropertyPart2(t *testing.T, ad Adapter) {

	// Create and fill original scheme instance
	origVal := reflect.New(reflect.TypeOf(ad).Elem())
	origScheme := origVal.Interface().(Adapter)
	fillDummyScheme(origScheme)

	// Render to string
	outStr, err := RenderAdapterToString(origScheme)
	if err != nil {
		t.Fatalf("RenderAdapterToString failed for %T: %v", origScheme, err)
	}

	// Create new empty scheme instance
	newVal := reflect.New(reflect.TypeOf(ad).Elem())
	newScheme := newVal.Interface().(Adapter)

	// Parse from string back into newScheme
	err = newScheme.FromString(outStr)
	if err != nil {
		t.Fatalf("FromString failed for %T: %v", newScheme, err)
	}

	// Check similarity using custom comparison that tolerates slight color differences
	f := FieldSimilarity(origScheme, newScheme)
	if f < 1.0 {
		t.Errorf("Round-trip scheme mismatch beyond tolerance")
		fmt.Printf("%f\n\n", f)
		fmt.Printf("\n\n")
		printSchemeContents("Original", origScheme)
		fmt.Printf("\n\n")
		printSchemeContents("Parsed", newScheme)
	}
}

type mockScheme struct {
	FieldA *string
	FieldB *string
	FieldC *string
}

func TestFieldSimilarity(t *testing.T) {
	a1 := "red"
	a2 := "blue"
	s1 := &mockScheme{FieldA: &a1, FieldB: &a2}
	s2 := &mockScheme{FieldA: &a1, FieldB: &a2}
	s3 := &mockScheme{FieldA: &a1, FieldB: nil}

	sim := FieldSimilarity(s1, s2)
	if sim != 1.0 {
		t.Errorf("expected 1.0 similarity, got %.2f", sim)
	}

	sim2 := FieldSimilarity(s1, s3)
	if sim2 >= 1.0 {
		t.Errorf("expected less than 1.0 similarity, got %.2f", sim2)
	}
}
