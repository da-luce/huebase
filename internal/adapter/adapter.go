package adapter

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strings"

	"github.com/da-luce/paletteport/internal/color"
	log "github.com/da-luce/paletteport/internal/logger"
	"github.com/da-luce/paletteport/internal/objectmap"
	"github.com/da-luce/paletteport/templates"
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
// Yes, it would be nice to enforce all Adapters be structs so that we can handle
// the mapping with the objectmap functionality (more control over missing src and dst fields)
// BUT, we can't call objectmap when doing that (complains about not knowing about type of thing)
// Also, goes against the grain of go... not idiomatic and would likely introduce future problems
type Adapter interface {
	Name() string                  // Return shorthand name
	FromString(input string) error // Parse input text into the adapter's struct
	TemplateName() string          // Return path to the output generation template file
}

// List of registered adapters
var Adapters = []Adapter{
	// &Base16Scheme{},
	&AlacrittyScheme{},
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

func onMissingField(fieldPath []string, srcVal reflect.Value) {
	pathStr := strings.Join(fieldPath, ".")
	log.Logger.Warn().
		Msgf(
			"Source field %q with value %v was unused during mapping",
			pathStr,
			srcVal.Interface(),
		)
}

func onUnusedIntoAbstract(fieldPath []string, srcVal reflect.Value) {
	pathStr := strings.Join(fieldPath, ".")
	log.Logger.Debug().
		Msgf("Abstract field %q with value %v was unused during mapping",
			pathStr,
			srcVal.Interface(),
		)
}

func onUnusedFromAbstract(fieldPath []string, srcVal reflect.Value) {
	pathStr := strings.Join(fieldPath, ".")
	log.Logger.Debug().
		Msgf("Abstract field %q with value %v was unused during mapping",
			pathStr,
			srcVal.Interface(),
		)
}

func onMissingDest(fieldPath []string, srcVal reflect.Value) {
	pathStr := strings.Join(fieldPath, ".")
	log.Logger.Warn().
		Msgf(
			"Destination field %q with value %v was unused during mapping",
			pathStr,
			srcVal.Interface(),
		)

}

func adaptScheme(reader Adapter, writer Adapter) error {
	// Convert reader adapter to abstract scheme
	var abstractTheme AbstractScheme
	if err := objectmap.MapInto(
		reader,
		&abstractTheme,
		onMissingField,
		onUnusedIntoAbstract,
		"abstract",
	); err != nil {
		return fmt.Errorf("failed to convert reader to abstract: %w", err)
	}

	// Fill in missing fields
	fillUnsetInGroups(&abstractTheme)

	if err := objectmap.MapFrom(
		&abstractTheme,
		writer,
		onUnusedFromAbstract,
		onMissingDest,
		"abstract",
	); err != nil {
		return fmt.Errorf("failed to convert abstract to writer: %w", err)
	}

	return nil
}

// Renders an Adapter to a string using its TemplatePath.
func RenderAdapterToString(a Adapter) (string, error) {
	templateFile := a.TemplateName()
	tmplData, err := templates.FS.ReadFile(templateFile)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("out").Parse(string(tmplData))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, a); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ConvertTheme[S Adapter, W Adapter](input string, reader S, writer W) (string, error) {
	if err := reader.FromString(input); err != nil {
		return "", fmt.Errorf("failed to parse input file: %w", err)
	}

	if err := adaptScheme(reader, writer); err != nil {
		return "", err
	}

	return RenderAdapterToString(writer)
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
