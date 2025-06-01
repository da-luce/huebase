package adapters

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdapters(t *testing.T) {
	// Define test cases with different schemes
	testCases := []struct {
		name     string
		filepath string
		scheme   Adapter
	}{
		{name: "Base16", filepath: "../../themes/base16.yaml", scheme: &Base16Scheme{}},
		{name: "Alacritty ", filepath: "../../themes/alacritty.toml", scheme: &AlacrittyScheme{}},
		{name: "Windows Terminal", filepath: "../../themes/wt.json", scheme: &WindowsTerminalScheme{}},
		{name: "Gogh", filepath: "../../themes/gogh.yml", scheme: &GoghScheme{}},
		{name: "iTerm", filepath: "../../themes/iterm.itermcolors", scheme: &ItermScheme{}},
	}

	// Run the tests for each scheme
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			t.Run("LoadTheme", func(t *testing.T) {
				err := loadTheme(test.scheme, test.filepath)
				if err != nil {
					t.Fatalf("Failed to load theme for %s: %v", test.name, err)
				}
			})

			t.Run("SaveTheme", func(t *testing.T) {
				buffer := &bytes.Buffer{}
				err := saveTheme(test.scheme, buffer)
				if err != nil {
					t.Fatalf("Failed to save theme for %s: %v", test.name, err)
				}
				if buffer.Len() == 0 {
					t.Fatalf("SaveTheme produced empty output for %s", test.name)
				}
			})

			t.Run("InversePropertyTest", func(t *testing.T) {
				inversePropertyTest(t, test.filepath, test.scheme)
			})

			t.Run("NonNoneFields", func(t *testing.T) {
				nonNoneFieldsTest(t, test.filepath, test.scheme)
			})

			t.Run("SimilarityTest", func(t *testing.T) {
				similarityTest(t, test.filepath, test.scheme)
			})
		})
	}
}

// Define a function to read the theme file and create an adapter instance
func loadTheme(adapter Adapter, path string) error {
	// Read the file content
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read theme file '%s': %w", path, err)
	}

	// Use the adapter's FromString method to parse the content
	err = adapter.FromString(string(data))
	if err != nil {
		return fmt.Errorf("failed to parse theme file '%s': %w", path, err)
	}

	return nil
}

// saveTheme tests the adapter's ToString method by writing to a fake buffer.
//
// Parameters:
//   - adapter: An implementation of the Adapter interface.
//   - buffer: A writable buffer to simulate saving the theme.
//
// Returns:
//   - An error if the ToString method fails or if the buffer cannot be written to.
func saveTheme(adapter Adapter, buffer *bytes.Buffer) error {
	// Ensure the adapter implements the necessary method
	if adapter == nil {
		return fmt.Errorf("adapter cannot be nil")
	}

	// Use the adapter's ToString method to get the theme content as a string
	data, err := adapter.ToString()
	if err != nil {
		return fmt.Errorf("failed to serialize theme: %w", err)
	}

	// Write the data to the buffer
	_, err = buffer.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write to buffer: %w", err)
	}

	return nil
}

// Generic test for Adapter instances
func inversePropertyTest(t *testing.T, filepath string, scheme Adapter) {

	// Load the theme into the adapter
	err := loadTheme(scheme, filepath)
	assert.NoError(t, err, "loadTheme should not produce an error")

	// Convert the scheme to AbstractScheme
	abstract, err := ToAbstract(scheme)
	assert.NoError(t, err, "ToAbstract should not produce an error")

	// Create a new instance of the same type as the input scheme
	newScheme := reflect.New(reflect.TypeOf(scheme).Elem()).Interface().(Adapter)

	// Convert back from AbstractScheme to the original scheme
	FromAbstract(&abstract, newScheme)

	// Ensure the original and reconstructed schemes are deeply equal
	assert.True(t, reflect.DeepEqual(scheme, newScheme),
		"The original scheme and the reconstructed scheme should match")
	CompareStructs(scheme, newScheme)
}

func nonNoneFieldsTest(t *testing.T, filepath string, scheme Adapter) {

	// Load the theme into the adapter
	err := loadTheme(scheme, filepath)
	assert.NoError(t, err, "loadTheme should not produce an error")

	// Convert the scheme to AbstractScheme
	abstract, err := ToAbstract(scheme)
	assert.NoError(t, err, "ToAbstract should not produce an error")

	// Count non-None fields
	nonNoneCount := countNonNullFields(abstract)

	// Assert that there are at least 8 non-None fields
	assert.GreaterOrEqual(t, nonNoneCount, 8, "Expected at least 8 non-None fields in the abstract theme")
}

// CompareStructs prints the differences between two structs, including field names
func CompareStructs(a, b interface{}) {
	// Ensure both are of the same type
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		fmt.Println("Error: Structs are of different types")
		return
	}

	// Get the value and type of the structs, handling pointers
	valA := reflect.ValueOf(a)
	valB := reflect.ValueOf(b)

	if valA.Kind() == reflect.Ptr {
		valA = valA.Elem()
	}
	if valB.Kind() == reflect.Ptr {
		valB = valB.Elem()
	}

	if valA.Kind() != reflect.Struct || valB.Kind() != reflect.Struct {
		fmt.Println("Error: Inputs must be structs or pointers to structs")
		return
	}

	typ := valA.Type()

	// Iterate through the fields
	for i := 0; i < valA.NumField(); i++ {
		fieldA := valA.Field(i)
		fieldB := valB.Field(i)
		fieldType := typ.Field(i) // Get field type
		fieldName := fieldType.Name

		// Compare field values
		if !reflect.DeepEqual(fieldA.Interface(), fieldB.Interface()) {
			fmt.Printf("Difference in field '%s' (type: %s):\n\tStruct A: %v\n\tStruct B: %v\n",
				fieldName, fieldType.Type, fieldA.Interface(), fieldB.Interface())
		}
	}
}

// Similarity Test
func similarityTest(t *testing.T, filepath string, scheme Adapter) {
	// Load the theme into the adapter
	err := loadTheme(scheme, filepath)
	assert.NoError(t, err, "loadTheme should not produce an error")

	// Convert to abstract representation
	abstract, err := ToAbstract(scheme)
	assert.NoError(t, err, "ToAbstract should not produce an error")

	// Create a new instance of the same type as the input scheme
	newScheme := reflect.New(reflect.TypeOf(scheme).Elem()).Interface().(Adapter)

	// Convert back from abstract representation
	FromAbstract(&abstract, newScheme)

	// Serialize both the original and reconstructed scheme
	originalData, err := scheme.ToString()
	assert.NoError(t, err, "ToString should not produce an error for the original scheme")

	reconstructedData, err := newScheme.ToString()
	assert.NoError(t, err, "ToString should not produce an error for the reconstructed scheme")

	print("Origional data:\n" + originalData)
	print("Reconrstructed data:\n" + reconstructedData)
	// Compute similarity percentage
	similarity := computeSimilarity(originalData, reconstructedData)

	// Define a threshold (e.g., 90% similarity)
	const similarityThreshold = 100.0
	assert.GreaterOrEqual(t, similarity, similarityThreshold,
		fmt.Sprintf("Similarity should be at least %.2f%%, but got %.2f%%", similarityThreshold, similarity))
}

// Compute the similarity between two strings as a percentage
func computeSimilarity(a, b string) float64 {
	// Normalize line endings
	a = strings.ReplaceAll(a, "\r\n", "\n")
	b = strings.ReplaceAll(b, "\r\n", "\n")

	// Compute the number of matching characters
	matches := 0
	minLen := min(len(a), len(b))

	for i := 0; i < minLen; i++ {
		if a[i] == b[i] {
			matches++
		}
	}

	// Use the longer length as the denominator
	maxLen := max(len(a), len(b))
	if maxLen == 0 {
		return 100.0 // If both are empty, they are fully similar
	}

	return (float64(matches) / float64(maxLen)) * 100.0
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
