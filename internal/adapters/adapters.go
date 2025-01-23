package adapters

import "github.com/da-luce/huebase/internal/color"

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
