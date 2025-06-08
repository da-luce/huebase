package adapter

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strings"

	"github.com/da-luce/paletteport/internal/adapter/alacritty"
	"github.com/da-luce/paletteport/internal/adapter/base16"
	"github.com/da-luce/paletteport/internal/adapter/gogh"
	"github.com/da-luce/paletteport/internal/adapter/iterm"
	"github.com/da-luce/paletteport/internal/adapter/windows_terminal"
	"github.com/da-luce/paletteport/internal/color"
	log "github.com/da-luce/paletteport/internal/logger"
	"github.com/da-luce/paletteport/internal/objectmap"
	"github.com/da-luce/paletteport/internal/structutil"
	"github.com/da-luce/paletteport/templates"
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
	&base16.Base16Scheme{},
	&alacritty.AlacrittyScheme{},
	&gogh.GoghScheme{},
	&iterm.ItermScheme{},
	&windows_terminal.WindowsTerminalScheme{},
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

	tmpl := template.New("out")

	// Add important functions
	tmpl.Funcs(template.FuncMap{
		"indent": templates.Indent,
	})

	tmpl, err = tmpl.Parse(string(tmplData))
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

func isColor(t reflect.Type) bool {
	// Dereference if pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name() == "Color" &&
		t.PkgPath() == "github.com/da-luce/paletteport/internal/color"
}

func isZeroValue(v reflect.Value) bool {
	// Invalid values are zero by definition
	if !v.IsValid() {
		return true
	}

	// If it's a pointer, check if nil
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		return v.IsNil()
	}

	// Use reflect.Zero to get zero value of the type and compare
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}

// FieldSimilarity compares two structs and returns the ratio of matching pointer field values,
// counting all fields where at least one side is set (non-nil).
// FIXME: the logic here sucks!!!
func FieldSimilarity(a, b any) float64 {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Kind() == reflect.Ptr {
		va = va.Elem()
	}
	if vb.Kind() == reflect.Ptr {
		vb = vb.Elem()
	}

	if va.Type() != vb.Type() {
		return 0.0
	}

	total := 0
	matching := 0

	structutil.TraverseStructDFS(a, func(path []string, field reflect.StructField, valA reflect.Value) bool {
		found, _, valB := structutil.HasNestedFieldSlice(vb, path)
		if !found {
			return true
		}

		if !valA.IsValid() && !valB.IsValid() {
			return true
		}

		// Dereference pointers once, if possible
		if valA.Kind() == reflect.Ptr && !valA.IsNil() {
			valA = valA.Elem()
		}
		if valB.Kind() == reflect.Ptr && !valB.IsNil() {
			valB = valB.Elem()
		}

		// Only handle base kinds or color pointers here
		kindA := valA.Kind()
		kindB := valB.Kind()

		// Helper to check if kind is a "base kind"
		isBaseKind := func(k reflect.Kind) bool {
			switch k {
			case reflect.String, reflect.Bool,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
				reflect.Float32, reflect.Float64:
				return true
			}
			return false
		}

		// Handle color pointers (need special check, as we deref above)
		if isColor(field.Type) {
			// valA and valB here are already dereferenced, so they should be color.Color structs, not pointers
			c1, ok1 := valA.Interface().(color.Color)
			c2, ok2 := valB.Interface().(color.Color)
			if !ok1 || !ok2 {
				// If either is not a color.Color value (e.g. zero), skip
				return true
			}

			total++
			if color.ColorsSimilar(c1, c2, 0.01) {
				matching++
			}
			return false
		}

		// If both sides are base kinds, compare them directly
		if isBaseKind(kindA) || isBaseKind(kindB) {
			nonZeroA := !isZeroValue(valA)
			nonZeroB := !isZeroValue(valB)
			if nonZeroA || nonZeroB {
				total++
				if reflect.DeepEqual(valA.Interface(), valB.Interface()) {
					matching++
				}
			}
			return false
		}

		// For any other kinds (struct, slice, map, etc), skip comparing (or count mismatch if desired)
		return true
	})

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
