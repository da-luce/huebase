package adapters

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestAdapters_ToAbstractAndFromAbstract(t *testing.T) {
	// Define test cases with different schemes
	testCases := []struct {
		name     string
		filepath string
		scheme   Adapter
	}{
		{
			name:     "Base16Scheme",
			filepath: "../../themes/base16.yaml",
			scheme:   &Base16Scheme{},
		},
		{
			name:     "AlacrittyScheme",
			filepath: "../../themes/alacritty.toml",
			scheme:   &AlacrittyScheme{},
		},
	}

	// Run the  test for each scheme
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			inversePropertyTest(t, test.filepath, test.scheme)
		})
	}
}

// Generic test for Adapter instances
func inversePropertyTest(t *testing.T, filepath string, scheme Adapter) {

	// Load the theme into the adapter
	err := loadTheme(scheme, filepath)
	assert.NoError(t, err, "loadTheme should not produce an error")

	// Convert the scheme to AbstractScheme
	abstract, err := ToAbstract(scheme)
	assert.NoError(t, err, "ToAbstract should not produce an error")

	t.Logf("Abstract Theme: %+v", abstract)

	// Create a new instance of the same type as the input scheme
	newScheme := reflect.New(reflect.TypeOf(scheme).Elem()).Interface().(Adapter)

	// Convert back from AbstractScheme to the original scheme
	FromAbstract(&abstract, newScheme)

	t.Logf("Original Scheme: %+v", scheme)
	t.Logf("Reconstructed Scheme: %+v", newScheme)

	// Ensure the original and reconstructed schemes are deeply equal
	assert.True(t, reflect.DeepEqual(scheme, newScheme),
		"The original scheme and the reconstructed scheme should match")
}
