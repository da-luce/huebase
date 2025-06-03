package adapter

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"reflect"

	"github.com/da-luce/huebase/internal/color"
	"github.com/da-luce/huebase/internal/fieldmap"
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
	FromString(input string) error // Parse input text into the adapter's struct
	MappingPath() string           // Return path to the YAML field mapping file
	TemplatePath() string          // Return path to the output generation template file
	Name() string                  // Return shorthand name
}

// List of registered adapters
var Adapters = []Adapter{
	&Base16Scheme{},
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
	Date   *string
}
type Scope struct {
	Basic         BasicScope
	Advanced      AdvancedScope
	Markup        MarkupScope
	Diagnostics   DiagnosticScope
	Editor        EditorScope
	Miscellaneous MiscScope
}

type BasicScope struct {
	Comment  *Color
	Keyword  *Color
	Constant *Color
	String   *Color
	Number   *Color
	Function *Color
	Variable *Color
	Operator *Color
}

type AdvancedScope struct {
	Class     *Color
	Type      *Color
	Property  *Color
	Attribute *Color
	Tag       *Color
	Namespace *Color
	Parameter *Color
	Selector  *Color
}

type MarkupScope struct {
	Heading     *Color
	Bold        *Color
	Italic      *Color
	Underline   *Color
	Link        *Color
	Quote       *Color
	List        *Color
	CodeBlock   *Color
	RawText     *Color
	TemplateTag *Color
}

type DiagnosticScope struct {
	Invalid    *Color
	Deprecated *Color
}

type EditorScope struct {
	Cursor      *Color
	CursorLine  *Color
	LineNumbers *Color
	Highlight   *Color
}

type MiscScope struct {
	Meta       *Color
	Annotation *Color
	Regex      *Color
	Background *Color
	Foreground *Color
}

// AbstractTheme represents the most generic theme possible
type AbstractScheme struct {
	Metadata      Meta
	ScopeColors   Scope
	AnsiColors    AnsiColors
	SpecialColors SpecialColors
}

func ConvertTheme(input string, reader Adapter, writer Adapter) (string, error) {
	// Parse input string into reader's struct
	if err := reader.FromString(input); err != nil {
		return "", fmt.Errorf("failed to parse input file: %w", err)
	}

	// Load the mapping from YAML string
	inputMapFile := reader.MappingPath()
	inputMapping, err := fieldmap.LoadMappingFromString(inputMapFile)
	if err != nil {
		return "", fmt.Errorf("failed to load mapping: %w", err)
	}

	// Prepare AbstractScheme to receive mapped values
	var abstractTheme AbstractScheme

	// Map reader -> abstract
	err = fieldmap.ApplySourceToDestMapping(reader, &abstractTheme, inputMapping)
	if err != nil {
		return "", fmt.Errorf("failed to map source to abstract: %w", err)
	}

	// Load the mapping from YAML string
	outputMapFile := writer.MappingPath()
	outputMapping, err := fieldmap.LoadMappingFromString(outputMapFile)
	if err != nil {
		return "", fmt.Errorf("failed to load mapping: %w", err)
	}

	// Map abstract -> writer
	err = fieldmap.ApplyDestToSourceMapping(&abstractTheme, writer, outputMapping)
	if err != nil {
		return "", fmt.Errorf("failed to map abstract to dest: %w", err)
	}

	// Load and parse the output template file
	tmplPath := writer.TemplatePath()
	tmplData, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", tmplPath, err)
	}

	tmpl, err := template.New("output").Parse(string(tmplData))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Render the template with writer data
	var outputBuf bytes.Buffer
	if err := tmpl.Execute(&outputBuf, writer); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return outputBuf.String(), nil
}

// FieldSimilarity compares two structs and returns the ratio of matching pointer field values,
// counting all fields where at least one side is set (non-nil).
func FieldSimilarity(a, b interface{}) float64 {
	va := reflect.ValueOf(a).Elem()
	vb := reflect.ValueOf(b).Elem()

	if va.Type() != vb.Type() {
		return 0.0
	}

	total := 0
	matching := 0

	for i := 0; i < va.NumField(); i++ {
		fa := va.Field(i)
		fb := vb.Field(i)

		if fa.Kind() != reflect.Ptr || fb.Kind() != reflect.Ptr {
			continue // skip non-pointer fields
		}

		// Count this field if either side is non-nil
		if !fa.IsNil() || !fb.IsNil() {
			total++

			// Count as match only if both are non-nil and deeply equal
			if !fa.IsNil() && !fb.IsNil() && reflect.DeepEqual(fa.Interface(), fb.Interface()) {
				matching++
			}
		}
	}

	if total == 0 {
		return 0.0
	}
	return float64(matching) / float64(total)
}

// // ToAbstract converts an Adapter interface implementation into an AbstractScheme.
// // Fields are mapped based on "abstract" tags in the concrete struct.
// //
// // Parameters:
// //   - input: An implementation of the Adapter interface.
// //
// // Returns:
// //   - AbstractScheme with mapped values.
// //   - An error if the input is invalid or mapping fails.
// func ToAbstract(adapter Adapter) (AbstractScheme, error) {

// 	// Ensure the adapter is a struct or pointer to a struct
// 	adapterValue := reflect.ValueOf(adapter)
// 	if adapterValue.Kind() == reflect.Ptr {
// 		adapterValue = adapterValue.Elem()
// 	}
// 	if adapterValue.Kind() != reflect.Struct {
// 		return AbstractScheme{}, fmt.Errorf("adapter must be a struct or pointer to a struct, got %s", adapterValue.Kind())
// 	}

// 	// Create a new AbstractScheme instance
// 	abstract := AbstractScheme{}
// 	abstractValue := reflect.ValueOf(&abstract).Elem()

// 	// Fill fields in abstract scheme
// 	traverseFields(adapterValue, func(field reflect.StructField, fieldValue reflect.Value) {

// 		// Get the 'abstract' tag
// 		tag := field.Tag.Get("abstract")
// 		if tag == "" {
// 			log.Printf("Field '%s' does not have an 'abstract' tag. Skipping.", field.Name)
// 			return
// 		}

// 		// Split the tag into parts (e.g., "AnsiColors.Red")
// 		tagParts := strings.Split(tag, ".")
// 		currentField := abstractValue

// 		// Traverse through nested fields in the AbstractScheme
// 		for _, part := range tagParts {
// 			currentField = currentField.FieldByName(part)
// 			if !abstractValue.IsValid() {
// 				log.Printf("Invalid field path '%s' in AbstractScheme for field '%s'.", tag, field.Name)
// 				return
// 			}
// 		}

// 		// We're at the last part of the tag, set the value
// 		if currentField.CanSet() {
// 			currentField.Set(fieldValue)
// 		} else {
// 			log.Printf("Cannot set value for '%s' in AbstractScheme. Field is not settable.", tag)
// 		}

// 	})

// 	return abstract, nil
// }

// // FromAbstract populates an Adapter implementation with values from an AbstractScheme.
// // Fields in the Adapter are set based on "abstract" tags in the concrete struct.
// //
// // Parameters:
// //   - abstract: A pointer to the AbstractScheme containing source data.
// //   - output: An implementation of the Adapter interface to be populated.
// func FromAbstract(abstract *AbstractScheme, output Adapter) error {
// 	// Ensure `abstract` is not nil
// 	if abstract == nil {
// 		return fmt.Errorf("Error: AbstractScheme is nil. Cannot populate fields")
// 	}

// 	outputValue := reflect.ValueOf(output)
// 	if outputValue.Kind() != reflect.Ptr || outputValue.Elem().Kind() != reflect.Struct {
// 		return fmt.Errorf("Output must be a pointer to a struct, got %s", outputValue.Kind())
// 	}

// 	// Get the underlying struct value
// 	outputValue = outputValue.Elem()

// 	// Use visitFields to iterate through all fields of the output struct
// 	traverseFields(outputValue, func(field reflect.StructField, fieldValue reflect.Value) {
// 		// Get the 'abstract' tag
// 		tag := field.Tag.Get("abstract")
// 		if tag == "" {
// 			log.Printf("Field '%s' does not have an 'abstract' tag. Skipping.", field.Name)
// 			return
// 		}

// 		// Split the tag into parts (e.g., "AnsiColors.Red")
// 		tagParts := strings.Split(tag, ".")
// 		abstractValue := reflect.ValueOf(abstract).Elem()

// 		// Traverse through nested fields in the AbstractScheme
// 		for _, part := range tagParts {
// 			abstractValue = abstractValue.FieldByName(part)
// 			if !abstractValue.IsValid() {
// 				log.Printf("Invalid field path '%s' in AbstractScheme for field '%s'.", tag, field.Name)
// 				return
// 			}
// 		}

// 		// Set value in the output struct if valid
// 		if abstractValue.IsValid() && abstractValue.CanInterface() {
// 			if fieldValue.CanSet() {
// 				fieldValue.Set(abstractValue)
// 			} else {
// 				log.Printf("Cannot set value for field '%s' in output struct.", field.Name)
// 			}
// 		}
// 	})

// 	return nil
// }

// // isBaseType determines if a reflect.Type is a base type that should not be further traversed.
// // Base types are used to set fields directly when mapping with "abstract" tags.
// //
// // Parameters:
// //   - t: The reflect.Type to check.
// //
// // Returns:
// //   - true if the type is a base type (e.g., int, string, etc.).
// //   - false if the type is a struct or another complex type that can be traversed.
// func isBaseType(t reflect.Type) bool {
// 	switch t.Kind() {
// 	case reflect.String, reflect.Int, reflect.Int64, reflect.Float64, reflect.Bool:
// 		return true
// 	case reflect.Ptr:
// 		// Check if the pointer points to a base type
// 		return isBaseType(t.Elem())
// 	default:
// 		if t.Name() == "Color" {
// 			return true
// 		}
// 		return false
// 	}
// }

// // traverseFields iterates over all fields in a struct or a pointer to a struct.
// // If the input is a pointer, it is dereferenced automatically.
// //
// // Parameters:
// //   - value: A reflect.Value of a struct or a pointer to a struct.
// //   - visitFunc: A function called for each field, receiving the field's metadata
// //     (reflect.StructField) and value (reflect.Value).
// func traverseFields(value reflect.Value, visitFunc func(field reflect.StructField, value reflect.Value)) {

// 	// Ensure the value is a struct or a pointer to a struct
// 	if value.Kind() == reflect.Ptr {
// 		value = value.Elem()
// 	}

// 	// Iterate over the fields of the struct
// 	for i := 0; i < value.NumField(); i++ {

// 		field := value.Type().Field(i)
// 		fieldValue := value.Field(i)

// 		// Visit the current field
// 		visitFunc(field, fieldValue)

// 		// Do not recurse further on Color structs and the like
// 		if isBaseType(field.Type) {
// 			continue
// 		}

// 		// Recursively visit nested structs
// 		if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Struct {
// 			traverseFields(fieldValue, visitFunc)
// 		}
// 	}
// }

// // countNonNullFields counts the number of non-None fields in an AbstractScheme.
// //
// // Parameters:
// //   - abstract: The AbstractScheme to analyze.
// //
// // Returns:
// //   - An integer representing the count of fields that are not None.
// func countNonNullFields(abstract AbstractScheme) int {
// 	count := 0

// 	traverseFields(reflect.ValueOf(abstract), func(field reflect.StructField, fieldValue reflect.Value) {
// 		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
// 			count++
// 		}
// 	})

// 	return count
// }

// // warnUnsetFields iterates through fields in a struct and logs a warning
// // for each field that is a pointer and has not been set (i.e., is nil).
// //
// // Parameters:
// //   - value: The reflect.Value of a struct or a pointer to a struct.
// //
// // Returns:
// //   - None.
// func WarnUnsetFields(value interface{}, prefix string, level LogLevel) {
// 	// Ensure that the input is a pointer to a struct
// 	adapterValue := reflect.ValueOf(value)
// 	if adapterValue.Kind() == reflect.Ptr {
// 		adapterValue = adapterValue.Elem()
// 	}

// 	if adapterValue.Kind() != reflect.Struct {
// 		log.Println("Error: Expected a struct, got", adapterValue.Kind())
// 		return
// 	}

// 	// Iterate through the fields and check if they are unset (nil)
// 	traverseFields(adapterValue, func(field reflect.StructField, fieldValue reflect.Value) {
// 		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
// 			switch level {
// 			case LevelInfo:
// 				log.Infof("%s: Field '%s' is not set.", prefix, field.Name)
// 			case LevelError:
// 				log.Errorf("%s: Field '%s' is not set.", prefix, field.Name)
// 			}
// 		}
// 	})
// }
