package adapter

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"reflect"

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
	Name() string // Return shorthand name
	toAbstract() (*AbstractScheme, error)
	fromAbstract(a *AbstractScheme) error
	FromString(input string) error // Parse input text into the adapter's struct
	TemplatePath() string          // Return path to the output generation template file
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

func adaptScheme(reader Adapter, writer Adapter) error {
	// Convert reader adapter to abstract scheme
	abstractTheme, err := reader.toAbstract()
	if err != nil {
		return fmt.Errorf("failed to convert reader to abstract: %w", err)
	}

	// Fill in missing fields
	fillUnsetInGroups(abstractTheme)

	// Convert abstract scheme back to writer adapter
	err = writer.fromAbstract(abstractTheme)
	if err != nil {
		return fmt.Errorf("failed to convert abstract to writer: %w", err)
	}

	return nil
}

func ConvertTheme[S Adapter, W Adapter](input string, reader S, writer W) (string, error) {
	// Parse input string into reader's struct
	if err := reader.FromString(input); err != nil {
		return "", fmt.Errorf("failed to parse input file: %w", err)
	}

	// Apply conversion
	adaptScheme(reader, writer)

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

// First pointer represents optionality of color, then we get a reference to the color
type ColorGroup []**color.Color

// TODO: expand these greatly
func fallbackGroups(s *AbstractScheme) []ColorGroup {
	return []ColorGroup{
		{
			&s.AnsiColors.Black,
			&s.SpecialColors.Background,
		},
		{
			&s.AnsiColors.Blue,
			&s.SpecialColors.Links,
		},
	}
}

func processGroups(s *AbstractScheme, groups []ColorGroup, groupFn func(ColorGroup)) {
	for _, group := range groups {
		groupFn(group)
	}
}

func fallbackGroup(group ColorGroup) {
	var fallback *color.Color
	for _, field := range group {
		if field != nil && *field != nil {
			fallback = *field
			break
		}
	}
	if fallback == nil {
		return // No fallback found
	}
	for _, field := range group {
		if field != nil && *field == nil {
			*field = fallback
		}
	}
}

func fillUnsetInGroups(s *AbstractScheme) {
	processGroups(s, fallbackGroups(s), fallbackGroup)
}
