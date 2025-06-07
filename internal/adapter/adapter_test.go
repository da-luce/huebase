package adapter

import (
	"fmt"
	"math/rand"
	"testing"

	"reflect"

	"github.com/da-luce/paletteport/internal/color"
)

func TestAllAdapters(t *testing.T) {
	for _, ad := range Adapters {
		t.Run(ad.Name(), func(t *testing.T) {
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

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func fillDummyScheme(a Adapter) {
	val := reflect.ValueOf(a).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			fieldType := field.Type().Elem()

			switch fieldType.Name() {
			case "Color":
				randomColor := color.RandomColor()
				ptr := reflect.New(fieldType)
				ptr.Elem().Set(reflect.ValueOf(randomColor))
				field.Set(ptr)
			case "string":
				randStr := randomString(8) // 8 char random string
				field.Set(reflect.ValueOf(&randStr))
			default:
				// For other struct pointers, create zero value pointer
				ptr := reflect.New(fieldType)
				field.Set(ptr)
			}
		}
	}
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
func printSchemeContents(prefix string, scheme interface{}) {
	val := reflect.ValueOf(scheme)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fmt.Printf("%s:\n", prefix)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			fmt.Printf("  %s: %v\n", fieldName, field.Elem())
		} else {
			fmt.Printf("  %s: %v\n", fieldName, field.Interface())
		}
	}
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
	if !AdaptersSimilar(origScheme, newScheme, 0.01) {
		t.Errorf("Round-trip scheme mismatch beyond tolerance")
		printSchemeContents("Original", origScheme)
		printSchemeContents("Parsed", newScheme)
	}
}

// AdaptersSimilar recursively compares two adapters field-by-field,
// using ColorsSimilar for Color fields and DeepEqual for others.
// FIXME: clean up this logic
func AdaptersSimilar(a1, a2 Adapter, tol float64) bool {
	v1 := reflect.ValueOf(a1).Elem()
	v2 := reflect.ValueOf(a2).Elem()

	if v1.Type() != v2.Type() {
		return false
	}

	for i := 0; i < v1.NumField(); i++ {
		f1 := v1.Field(i)
		f2 := v2.Field(i)

		// Handle nil pointers
		if f1.Kind() == reflect.Ptr && f2.Kind() == reflect.Ptr {
			if f1.IsNil() && f2.IsNil() {
				continue
			}
			if f1.IsNil() || f2.IsNil() {
				return false
			}
		}

		// Check if field is *color.Color (or your specific Color type)
		if f1.Type().Elem().Name() == "Color" && f1.Kind() == reflect.Ptr {
			c1 := f1.Elem().Interface().(color.Color)
			c2 := f2.Elem().Interface().(color.Color)
			if !color.ColorsSimilar(c1, c2, tol) {
				return false
			}
			continue
		}

		// For other pointer fields, recursively compare underlying values if structs
		if f1.Kind() == reflect.Ptr && f1.Elem().Kind() == reflect.Struct {
			if !AdaptersSimilar(f1.Interface().(Adapter), f2.Interface().(Adapter), tol) {
				return false
			}
			continue
		}

		// For other fields, use DeepEqual
		if !reflect.DeepEqual(f1.Interface(), f2.Interface()) {
			return false
		}
	}

	return true
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
