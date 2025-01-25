package adapters

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
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
		{
			name:     "Base16",
			filepath: "../../themes/base16.yaml",
			scheme:   &Base16Scheme{},
		},
		{
			name:     "Alacritty ",
			filepath: "../../themes/alacritty.toml",
			scheme:   &AlacrittyScheme{},
		},
		{
			name:     "Windows Terminal",
			filepath: "../../themes/wt.json",
			scheme:   &WindowsTerminalScheme{},
		},
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
}

func nonNoneFieldsTest(t *testing.T, filepath string, scheme Adapter) {

	// Load the theme into the adapter
	err := loadTheme(scheme, filepath)
	assert.NoError(t, err, "loadTheme should not produce an error")

	// Convert the scheme to AbstractScheme
	abstract, err := ToAbstract(scheme)
	assert.NoError(t, err, "ToAbstract should not produce an error")

	// Count non-None fields
	nonNoneCount := CountNonNoneFields(abstract)

	// Assert that there are at least 8 non-None fields
	assert.GreaterOrEqual(t, nonNoneCount, 8, "Expected at least 8 non-None fields in the abstract theme")
}
