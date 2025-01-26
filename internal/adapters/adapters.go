package adapters

import (
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/da-luce/huebase/internal/color"
)

type LogLevel = int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Scheme adapter structs must implement this interface for reading and writing
// to text
type Adapter interface {
	FromString(input string) error
	ToString() (string, error)
}

type Color = color.Color
type AnsiColors struct {
	Black         *Color
	Red           *Color
	Green         *Color
	Yellow        *Color
	Blue          *Color
	Magenta       *Color
	Cyan          *Color
	White         *Color
	BrightBlack   *Color
	BrightRed     *Color
	BrightGreen   *Color
	BrightYellow  *Color
	BrightBlue    *Color
	BrightMagenta *Color
	BrightCyan    *Color
	BrightWhite   *Color
}

type SpecialColors struct {
	Foreground       *Color
	ForegroundBright *Color
	Background       *Color
	Cursor           *Color
	CursorText       *Color
	Selection        *Color
	SelectedText     *Color
	Links            *Color
	FindMatch        *Color
}

type Meta struct {
	Name   *string
	Author *string
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
func FromAbstract(abstract *AbstractScheme, output Adapter) error {
	// Ensure `abstract` is not nil
	if abstract == nil {
		return fmt.Errorf("Error: AbstractScheme is nil. Cannot populate fields")
	}

	outputValue := reflect.ValueOf(output)
	if outputValue.Kind() != reflect.Ptr || outputValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Output must be a pointer to a struct, got %s", outputValue.Kind())
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

	return nil
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

	// Iterate over the fields of the struct
	for i := 0; i < value.NumField(); i++ {

		field := value.Type().Field(i)
		fieldValue := value.Field(i)

		// Visit the current field
		visitFunc(field, fieldValue)

		// Do not recurse further on Color structs and the like
		if isBaseType(field.Type) {
			continue
		}

		// Recursively visit nested structs
		if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Struct {
			traverseFields(fieldValue, visitFunc)
		}
	}
}

// countNonNullFields counts the number of non-None fields in an AbstractScheme.
//
// Parameters:
//   - abstract: The AbstractScheme to analyze.
//
// Returns:
//   - An integer representing the count of fields that are not None.
func countNonNullFields(abstract AbstractScheme) int {
	count := 0

	traverseFields(reflect.ValueOf(abstract), func(field reflect.StructField, fieldValue reflect.Value) {
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			count++
		}
	})

	return count
}

// warnUnsetFields iterates through fields in a struct and logs a warning
// for each field that is a pointer and has not been set (i.e., is nil).
//
// Parameters:
//   - value: The reflect.Value of a struct or a pointer to a struct.
//
// Returns:
//   - None.
func WarnUnsetFields(value interface{}, prefix string, level LogLevel) {
	// Ensure that the input is a pointer to a struct
	adapterValue := reflect.ValueOf(value)
	if adapterValue.Kind() == reflect.Ptr {
		adapterValue = adapterValue.Elem()
	}

	if adapterValue.Kind() != reflect.Struct {
		log.Println("Error: Expected a struct, got", adapterValue.Kind())
		return
	}

	// Iterate through the fields and check if they are unset (nil)
	traverseFields(adapterValue, func(field reflect.StructField, fieldValue reflect.Value) {
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			switch level {
			case LevelInfo:
				log.Infof("%s: Field '%s' is not set.", prefix, field.Name)
			case LevelError:
				log.Errorf("%s: Field '%s' is not set.", prefix, field.Name)
			}
		}
	})
}
