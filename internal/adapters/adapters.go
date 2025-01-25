package adapters

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/da-luce/huebase/internal/color"
	"github.com/da-luce/huebase/internal/option"
)

// Scheme adapter structs must implement this interface for reading and writing
// to text
type Adapter interface {
	ToString() (string, error)
	FromString(input string) error
}

type Color = color.Color
type OptColor = option.Option[Color]
type OptString = option.Option[string]
type AnsiColors struct {
	Black         OptColor
	Red           OptColor
	Green         OptColor
	Yellow        OptColor
	Blue          OptColor
	Magenta       OptColor
	Cyan          OptColor
	White         OptColor
	BrightBlack   OptColor
	BrightRed     OptColor
	BrightGreen   OptColor
	BrightYellow  OptColor
	BrightBlue    OptColor
	BrightMagenta OptColor
	BrightCyan    OptColor
	BrightWhite   OptColor
}

type SpecialColors struct {
	Foreground       OptColor
	ForegroundBright OptColor
	Background       OptColor
	Cursor           OptColor
	CursorText       OptColor
	Selection        OptColor
	SelectedText     OptColor
	Links            OptColor
	FindMatch        OptColor
}

type Meta struct {
	Name   OptString
	Author OptString
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
func ToAbstract(adapter Adapter) (AbstractScheme, error) {

	// Ensure the adapter is a struct or pointer to a struct
	adapterValue := reflect.ValueOf(adapter)
	if adapterValue.Kind() == reflect.Ptr {
		adapterValue = adapterValue.Elem()
	}
	if adapterValue.Kind() != reflect.Struct {
		return AbstractScheme{}, fmt.Errorf("adapter must be a struct or pointer to a struct, got %s", adapterValue.Kind())
	}

	// Create a new AbstractScheme instance
	abstract := AbstractScheme{}
	abstractValue := reflect.ValueOf(&abstract).Elem()

	// Fill fields in abstract scheme
	traverseFields(adapterValue, func(field reflect.StructField, fieldValue reflect.Value) {

		// Get the 'abstract' tag
		tag := field.Tag.Get("abstract")
		if tag == "" {
			log.Printf("Field '%s' does not have an 'abstract' tag. Skipping.", field.Name)
			return
		}

		// Split the tag into parts (e.g., "AnsiColors.Red")
		tagParts := strings.Split(tag, ".")
		currentField := abstractValue

		// Traverse through nested fields in the AbstractScheme
		for _, part := range tagParts {
			currentField = currentField.FieldByName(part)
			if !abstractValue.IsValid() {
				log.Printf("Invalid field path '%s' in AbstractScheme for field '%s'.", tag, field.Name)
				return
			}
		}

		// We're at the last part of the tag, set the value
		if currentField.CanSet() {
			currentField.Set(fieldValue)
		} else {
			log.Printf("Cannot set value for '%s' in AbstractScheme. Field is not settable.", tag)
		}

	})

	return abstract, nil
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
		if strings.HasPrefix(t.Name(), "Option") {
			return true
		}
		return false
	}
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
		if isBaseType(field.Type) || isOption(field) {
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

func isOption(x interface{}) bool {
	// FIXME: this is gross
	if x == nil {
		return false
	}

	typ := reflect.TypeOf(x)
	return strings.HasPrefix(typ.Name(), "Option")
}

// CountNonNoneFields counts the number of non-None fields in an AbstractScheme.
//
// Parameters:
//   - abstract: The AbstractScheme to analyze.
//
// Returns:
//   - An integer representing the count of fields that are not None.
func CountNonNoneFields(abstract AbstractScheme) int {
	count := 0

	// Traverse all fields in the AbstractScheme
	traverseFields(reflect.ValueOf(abstract), func(field reflect.StructField, fieldValue reflect.Value) {
		// Check if the field is an Option type
		if isOption(fieldValue.Interface()) {
			// Check if the isSet field in the Option struct is true
			if fieldValue.FieldByName("isSet").Bool() {
				count++
			}
		}
	})

	return count
}
