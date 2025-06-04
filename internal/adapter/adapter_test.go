package adapter

import (
	"math/rand"
	"testing"

	"reflect"

	"github.com/da-luce/huebase/internal/color"
)

func TestAllAdapters(t *testing.T) {
	for _, ad := range Adapters {
		t.Run(ad.Name(), func(t *testing.T) {
			t.Run("TransitiveProperty", func(t *testing.T) {
				testTransitiveProperty(t, ad, 0.9)
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
