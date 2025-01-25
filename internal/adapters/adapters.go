package adapters

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/da-luce/huebase/internal/color"
)

type Adapter interface {
	ToString() (string, error)
	FromString(input string) error
}

// TODO: add nil or Color
// AnsiColors represents the ANSI color palette.
type Color = color.Color
type AnsiColors struct {
	Black         Color
	Red           Color
	Green         Color
	Yellow        Color
	Blue          Color
	Magenta       Color
	Cyan          Color
	White         Color
	BrightBlack   Color
	BrightRed     Color
	BrightGreen   Color
	BrightYellow  Color
	BrightBlue    Color
	BrightMagenta Color
	BrightCyan    Color
	BrightWhite   Color
}

// SpecialColors represents special colors used in themes.
type SpecialColors struct {
	Foreground       Color
	ForegroundBright Color
	Background       Color
	Cursor           Color
	CursorText       Color
	Selection        Color
	SelectedText     Color
	Links            Color
	FindMatch        Color
}

// Meta contains metadata about the theme.
type Meta struct {
	Name   string
	Author string
}

// AbstractTheme represents the most generic theme possible
type AbstractScheme struct {
	AnsiColors    AnsiColors
	SpecialColors SpecialColors
	Metadata      Meta
}

// ToAbstract converts an Adapter interface implementation into an AbstractScheme.
// Fields are mapped based on "abstract" tags in the concrete struct.
//
// Parameters:
//   - input: An implementation of the Adapter interface.
//
// Returns:
//   - AbstractScheme with mapped values.
//   - An error if the input is invalid or mapping fails.
func ToAbstract(input Adapter) (AbstractScheme, error) {
	// Use reflection to inspect the input struct
	inputValue := reflect.ValueOf(input)

	// If the input is a pointer, dereference it
	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}

	// Ensure input is now a struct
	if inputValue.Kind() != reflect.Struct {
		return AbstractScheme{}, fmt.Errorf("input must be a struct or a pointer to a struct, got %s", inputValue.Kind())
	}

	// Create a new AbstractScheme instance
	abstract := AbstractScheme{}
	abstractValue := reflect.ValueOf(&abstract).Elem() // Dereference the pointer to AbstractScheme

	// Use visitFields to iterate over input fields
	traverseFields(inputValue, func(field reflect.StructField, fieldValue reflect.Value) {
		// Get the 'abstract' tag
		tag := field.Tag.Get("abstract")
		if tag == "" {
			log.Printf("Field '%s' does not have an 'abstract' tag. Skipping.", field.Name)
			return
		}

		// Split the tag into parts (e.g., "AnsiColors.Red")
		tagParts := strings.Split(tag, ".")
		currentField := abstractValue

		// Traverse through nested fields in AbstractScheme
		for i, part := range tagParts {
			currentField = currentField.FieldByName(part)

			// Handle invalid field paths
			if !currentField.IsValid() {
				log.Printf("Invalid field path '%s' in AbstractScheme for input field '%s' at part '%s'. Skipping.", tag, field.Name, part)
				return
			}

			// If we're at the last part of the tag, set the value
			if i == len(tagParts)-1 {
				if currentField.CanSet() {
					currentField.Set(fieldValue)
				} else {
					log.Printf("Cannot set value for '%s' in AbstractScheme. Field is not settable.", tag)
				}
			}
		}
	})

	return abstract, nil
}

// traverseFields iterates over all fields in a struct or a pointer to a struct.
// If the input is a pointer, it is dereferenced automatically.
//
// Parameters:
//   - value: A reflect.Value of a struct or a pointer to a struct.
//   - visitFunc: A function called for each field, receiving the field's metadata
//     (reflect.StructField) and value (reflect.Value).
func traverseFields(value reflect.Value, visitFunc func(field reflect.StructField, value reflect.Value)) {

	// Ensure the value is a struct or a pointer to a struct
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return
	}

	// Iterate over the fields of the struct
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i)

		// Do not recurse further on Color structs and the like
		if isBaseType(field.Type) {
			visitFunc(field, fieldValue)
			continue
		}

		// Visit the current field
		visitFunc(field, fieldValue)

		// Recursively visit nested structs
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() && fieldValue.Elem().Kind() == reflect.Struct {
			traverseFields(fieldValue, visitFunc)
		} else if fieldValue.Kind() == reflect.Struct {
			traverseFields(fieldValue, visitFunc)
		}
	}
}

// isBaseType determines if a reflect.Type is a base type that should not be further traversed.
// Base types are used to set fields directly when mapping with "abstract" tags.
//
// Parameters:
//   - t: The reflect.Type to check.
//
// Returns:
//   - true if the type is a base type (e.g., int, string, etc.).
//   - false if the type is a struct or another complex type that can be traversed.
func isBaseType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String, reflect.Int, reflect.Int64, reflect.Float64, reflect.Bool:
		return true
	case reflect.Ptr:
		// Check if the pointer points to a base type
		return isBaseType(t.Elem())
	default:
		if t.Name() == "Color" {
			return true
		}
		return false
	}
}

// FromAbstract populates an Adapter implementation with values from an AbstractScheme.
// Fields in the Adapter are set based on "abstract" tags in the concrete struct.
//
// Parameters:
//   - abstract: A pointer to the AbstractScheme containing source data.
//   - output: An implementation of the Adapter interface to be populated.
func FromAbstract(abstract *AbstractScheme, output Adapter) {
	// Ensure `abstract` is not nil
	if abstract == nil {
		log.Println("Error: AbstractScheme is nil. Cannot populate fields.")
		return
	}

	outputValue := reflect.ValueOf(output)
	if outputValue.Kind() != reflect.Ptr || outputValue.Elem().Kind() != reflect.Struct {
		log.Fatalf("Output must be a pointer to a struct, got %s", outputValue.Kind())
	}

	// Get the underlying struct value
	outputValue = outputValue.Elem()

	// Use visitFields to iterate through all fields of the output struct
	traverseFields(outputValue, func(field reflect.StructField, fieldValue reflect.Value) {
		// Get the 'abstract' tag
		tag := field.Tag.Get("abstract")
		if tag == "" {
			log.Printf("Field '%s' does not have an 'abstract' tag. Skipping.", field.Name)
			return
		}

		// Split the tag into parts (e.g., "AnsiColors.Red")
		tagParts := strings.Split(tag, ".")
		abstractValue := reflect.ValueOf(abstract).Elem()

		// Traverse through nested fields in the AbstractScheme
		for _, part := range tagParts {
			abstractValue = abstractValue.FieldByName(part)
			if !abstractValue.IsValid() {
				log.Printf("Invalid field path '%s' in AbstractScheme for field '%s'.", tag, field.Name)
				return
			}
		}

		// Set value in the output struct if valid
		if abstractValue.IsValid() && abstractValue.CanInterface() {
			if fieldValue.CanSet() {
				fieldValue.Set(abstractValue)
			} else {
				log.Printf("Cannot set value for field '%s' in output struct.", field.Name)
			}
		}
	})
}
