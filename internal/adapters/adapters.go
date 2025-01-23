package adapters

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/da-luce/huebase/internal/color"
)

type Adapter interface {
	ToString(input interface{}) (string, error)
	FromString(input string) (Adapter, error)
}

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
	Foreground   Color
	Background   Color
	Cursor       Color
	CursorText   Color
	Selection    Color
	SelectedText Color
	Links        Color
	FindMatch    Color
}

// Meta contains metadata about the theme.
type Meta struct {
	Name   string
	Author string
}

// AbstractTheme represents a theme with ANSI colors, special colors, and metadata.
type AbstractScheme struct {
	AnsiColors    AnsiColors
	SpecialColors SpecialColors
	Metadata      Meta
}

func ToAbstract(input interface{}) (*AbstractScheme, error) {
	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()

	// Ensure input is a struct
	if inputType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct, got %s", inputType.Kind())
	}

	// Create a new AbstractScheme and use a pointer to it
	abstract := &AbstractScheme{}
	abstractValue := reflect.ValueOf(abstract).Elem()

	// Iterate through fields in the input struct
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		tag := field.Tag.Get("abstract") // Look for the 'abstract' tag

		// Skip this field if it doesn't have an "abstract" tag
		if tag == "" {
			log.Printf("Field '%s' does not have an 'abstract' tag in conversion to abstract theme. Skipping.", field.Name)
			continue
		}

		// Split the tag into parts (e.g., "AnsiColors.Red")
		tagParts := strings.Split(tag, ".")
		abstractField := abstractValue

		// Traverse through nested fields
		for _, part := range tagParts {
			abstractField = abstractField.FieldByName(part)
			if !abstractField.IsValid() {
				log.Printf("Invalid field path '%s' in AbstractScheme for input field '%s'", tag, field.Name)
				break
			}
		}

		// Set value if the field path is valid
		if abstractField.IsValid() && abstractField.CanSet() {
			abstractField.Set(inputValue.Field(i))
		} else if abstractField.IsValid() {
			log.Printf("Cannot set value for '%s' in AbstractScheme", tag)
		}
	}

	return abstract, nil
}

func FromAbstract(abstract *AbstractScheme, output interface{}) {
	// Ensure `abstract` is not nil
	if abstract == nil {
		log.Println("Error: AbstractScheme is nil. Cannot populate fields.")
		return
	}

	outputValue := reflect.ValueOf(output)
	if outputValue.Kind() != reflect.Ptr || outputValue.Elem().Kind() != reflect.Struct {
		log.Fatalf("Output must be a pointer to a struct, got %s", outputValue.Kind())
	}

	outputValue = outputValue.Elem()
	outputType := outputValue.Type()

	// Iterate through fields in the output struct
	for i := 0; i < outputType.NumField(); i++ {
		field := outputType.Field(i)
		tag := field.Tag.Get("abstract") // Look for the 'abstract' tag

		// Skip this field if it doesn't have an "abstract" tag
		if tag == "" {
			log.Printf("Field '%s' does not have an 'abstract' tag in conversion to theme. Skipping.", field.Name)
			continue
		}

		// Split the tag into parts (e.g., "AnsiColors.Red")
		tagParts := strings.Split(tag, ".")
		abstractField := reflect.ValueOf(abstract).Elem()

		// Traverse through nested fields
		for _, part := range tagParts {
			abstractField = abstractField.FieldByName(part)
			if !abstractField.IsValid() {
				log.Printf("Invalid field path '%s' in AbstractScheme for output field '%s'", tag, field.Name)
				break
			}
		}

		// Set value in the output struct if the field path is valid
		if abstractField.IsValid() && abstractField.CanInterface() {
			outputField := outputValue.FieldByName(field.Name)
			if outputField.IsValid() && outputField.CanSet() {
				outputField.Set(abstractField)
			} else {
				log.Printf("Cannot set value for field '%s' in output struct", field.Name)
			}
		}
	}
}
