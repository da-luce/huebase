package adapter

import (
	"os"
	"testing"

	"reflect"

	"github.com/da-luce/huebase/internal/fieldmap"
)

func TestVerifyAllSchemeMappings(t *testing.T) {
	for _, adapter := range Adapters {
		data, err := os.ReadFile(adapter.MappingPath())
		if err != nil {
			t.Errorf("Failed to read %s: %v", adapter.MappingPath(), err)
			continue
		}
		err = fieldmap.VerifyMappingString(
			reflect.TypeOf(adapter).Elem(),
			reflect.TypeOf(AbstractScheme{}),
			string(data),
		)
		if err != nil {
			t.Errorf("Verification failed for %s: %v", adapter.MappingPath(), err)
		}
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

func fillDummyScheme(a Adapter) {
	val := reflect.ValueOf(a).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			// Assume pointer to string for simplicity here, replace with correct type if needed
			// For *Color, you'd create a dummy Color struct pointer instead
			fieldType := field.Type().Elem()

			switch fieldType.Kind() {
			case reflect.String:
				dummy := "dummy"
				field.Set(reflect.ValueOf(&dummy))
			case reflect.Struct:
				// If *Color or any struct pointer
				ptr := reflect.New(fieldType)
				field.Set(ptr)
			// add more cases as needed
			default:
				// Skip or add default initialization here if you want
			}
		}
	}
}

func TestTransitiveProperty(t *testing.T) {
	const similarityThreshold = 0.9 // 90%

	for _, adapter := range Adapters {
		// Create new instance of adapter type
		schemeVal := reflect.New(reflect.TypeOf(adapter).Elem())
		scheme := schemeVal.Interface().(Adapter)

		// Fill dummy data so all pointer fields are non-nil
		fillDummyScheme(scheme)

		// Load mapping file
		mappingData, err := os.ReadFile(scheme.MappingPath())
		if err != nil {
			t.Fatalf("Failed to read mapping for %T: %v", scheme, err)
		}
		mapping, err := fieldmap.LoadMappingFromString(string(mappingData))
		if err != nil {
			t.Fatalf("Failed to load mapping for %T: %v", scheme, err)
		}

		// Convert scheme -> AbstractScheme
		abstract := &AbstractScheme{}
		err = fieldmap.ApplySourceToDestMapping(scheme, abstract, mapping)
		if err != nil {
			t.Fatalf("Failed to map source to abstract for %T: %v", scheme, err)
		}

		// Convert AbstractScheme -> scheme
		newSchemeVal := reflect.New(reflect.TypeOf(adapter).Elem())
		newScheme := newSchemeVal.Interface().(Adapter)
		err = fieldmap.ApplyDestToSourceMapping(abstract, newScheme, mapping)
		if err != nil {
			t.Fatalf("Failed to map abstract back to source for %T: %v", scheme, err)
		}

		// Check that all pointer fields in newScheme are non-nil
		checkAllFieldsSet(t, newScheme)

		// Check similarity between original and new scheme
		sim := FieldSimilarity(scheme, newScheme)
		if sim < similarityThreshold {
			t.Errorf("Similarity too low for %T: got %.2f, want >= %.2f", scheme, sim, similarityThreshold)
		}
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
