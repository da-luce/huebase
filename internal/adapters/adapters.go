package adapters

import (
	"log"
	"reflect"
	"strings"

	"github.com/da-luce/huebase/internal/color"
)

// Reader interface for reading themes
type Reader interface {
	FromString(input string) (AbstractScheme, error)
}

// Writer interface for writing themes
type Writer interface {
	ToString(theme AbstractScheme) (string, error)
}

type Base16Adapter struct{}
type AlacrittyAdapter struct{}

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

func PopulateAbstractScheme(input interface{}, abstract *AbstractScheme) {
	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()

	if inputType.Kind() != reflect.Struct {
		log.Fatalf("Input must be a struct, got %s", inputType.Kind())
	}

	abstractValue := reflect.ValueOf(abstract).Elem()

	// Iterate through fields in the input struct
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		tag := field.Tag.Get("abstract") // Look for the 'abstract' tag

		// Skip this field if it doesn't have an "abstract" tag
		if tag == "" {
			log.Printf("Field '%s' does not have an 'abstract' tag. Skipping.", field.Name)
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
}
