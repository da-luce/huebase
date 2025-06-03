package adapter

import (
	"math/rand"
	"os"
	"testing"

	"reflect"

	"github.com/da-luce/huebase/internal/color"
	"github.com/da-luce/huebase/internal/fieldmap"
)

func TestAllAdapters(t *testing.T) {
	for _, ad := range Adapters {
		t.Run(ad.Name(), func(t *testing.T) {
			t.Run("VerifyAllSchemeMappings", func(t *testing.T) {
				testVerifyAllSchemeMappings(t, ad)
			})
			t.Run("AllConcreteFieldsAreMapped", func(t *testing.T) {
				testAllConcreteFieldsAreMapped(t, ad)
			})
			t.Run("AllAbstractFieldsAreCoveredInMapping", func(t *testing.T) {
				testAllAbstractFieldsAreCoveredInMapping(t, ad)
			})
			t.Run("TransitiveProperty", func(t *testing.T) {
				testTransitiveProperty(t, ad, 0.9)
			})
		})
	}
}

func testVerifyAllSchemeMappings(t *testing.T, ad Adapter) {
	data, err := os.ReadFile(ad.MappingPath())
	if err != nil {
		t.Errorf("Failed to read %s: %v", ad.MappingPath(), err)
	}
	err = fieldmap.VerifyMappingString(
		reflect.TypeOf(ad).Elem(),
		reflect.TypeOf(AbstractScheme{}),
		string(data),
	)
	if err != nil {
		t.Errorf("Verification failed for %s: %v", ad.MappingPath(), err)
	}
}

func testAllConcreteFieldsAreMapped(t *testing.T, ad Adapter) {
	data, err := os.ReadFile(ad.MappingPath())
	if err != nil {
		t.Errorf("Failed to read mapping file for adapter %s: %v", ad.Name(), err)
	}

	mapping, err := fieldmap.LoadMappingFromString(string(data))
	if err != nil {
		t.Errorf("Failed to load mapping for adapter %s: %v", ad.Name(), err)
	}

	mappedFields := map[string]bool{}
	for key := range mapping {
		mappedFields[key] = true
	}

	schemeType := reflect.TypeOf(ad).Elem()
	for i := 0; i < schemeType.NumField(); i++ {
		fieldName := schemeType.Field(i).Name
		if _, found := mappedFields[fieldName]; !found {
			t.Errorf("Field %s in %s is not represented in the mapping keys", fieldName, ad.Name())
		}
	}
}

func testAllAbstractFieldsAreCoveredInMapping(t *testing.T, ad Adapter) {
	data, err := os.ReadFile(ad.MappingPath())
	if err != nil {
		t.Errorf("Failed to read mapping file for adapter %s: %v", ad.Name(), err)
	}

	mapping, err := fieldmap.LoadMappingFromString(string(data))
	if err != nil {
		t.Errorf("Failed to load mapping for adapter %s: %v", ad.Name(), err)
	}

	// Collect all used abstract fields
	used := map[string]bool{}
	for _, abstractFields := range mapping {
		for _, name := range abstractFields {
			used[name] = true
		}
	}

	absType := reflect.TypeOf(AbstractScheme{})
	for i := 0; i < absType.NumField(); i++ {
		absField := absType.Field(i).Name
		if !used[absField] {
			t.Errorf("Abstract field %s not covered in mapping values for adapter %s", absField, ad.Name())
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
	newSchemeVal := reflect.New(reflect.TypeOf(ad).Elem())
	newScheme := newSchemeVal.Interface().(Adapter)
	err = fieldmap.ApplyDestToSourceMapping(abstract, newScheme, mapping)
	if err != nil {
		t.Fatalf("Failed to map abstract back to source for %T: %v", scheme, err)
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
