package adapters

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

// Primary colors
type Primary struct {
	Background Color `toml:"background" abstract:"SpecialColors.Background"`
	Foreground Color `toml:"foreground" abstract:"SpecialColors.Foreground"`
}

// Cursor colors
type Cursor struct {
	Cursor Color `toml:"cursor" abstract:"SpecialColors.Cursor"`
	Text   Color `toml:"text" abstract:"SpecialColors.CursorText"`
}

// Normal ANSI colors
type Normal struct {
	Black   Color `toml:"black" abstract:"AnsiColors.Black"`
	Blue    Color `toml:"blue" abstract:"AnsiColors.Blue"`
	Cyan    Color `toml:"cyan" abstract:"AnsiColors.Cyan"`
	Green   Color `toml:"green" abstract:"AnsiColors.Green"`
	Magenta Color `toml:"magenta" abstract:"AnsiColors.Magenta"`
	Red     Color `toml:"red" abstract:"AnsiColors.Red"`
	White   Color `toml:"white" abstract:"AnsiColors.White"`
	Yellow  Color `toml:"yellow" abstract:"AnsiColors.Yellow"`
}

// Bright ANSI colors
type Bright struct {
	Black   Color `toml:"black" abstract:"AnsiColors.BrightBlack"`
	Blue    Color `toml:"blue" abstract:"AnsiColors.BrightBlue"`
	Cyan    Color `toml:"cyan" abstract:"AnsiColors.BrightCyan"`
	Green   Color `toml:"green" abstract:"AnsiColors.BrightGreen"`
	Magenta Color `toml:"magenta" abstract:"AnsiColors.BrightMagenta"`
	Red     Color `toml:"red" abstract:"AnsiColors.BrightRed"`
	White   Color `toml:"white" abstract:"AnsiColors.BrightWhite"`
	Yellow  Color `toml:"yellow" abstract:"AnsiColors.BrightYellow"`
}

// Selection colors
type Selection struct {
	Background Color `toml:"background" abstract:"SpecialColors.Selection"`
	Text       Color `toml:"text" abstract:"SpecialColors.SelectedText"`
}

// Colors struct
type Colors struct {
	Primary   Primary   `toml:"primary"`
	Cursor    Cursor    `toml:"cursor"`
	Normal    Normal    `toml:"normal"`
	Bright    Bright    `toml:"colors.bright"`
	Selection Selection `toml:"selection"`
}

// AlacrittyScheme with abstract tags
type AlacrittyScheme struct {
	Colors Colors `toml:"colors"`
}

func (a *AlacrittyScheme) FromString(input string) error {
	// Unmarshal the input TOML string into the AlacrittyScheme struct
	if err := toml.Unmarshal([]byte(input), a); err != nil {
		return fmt.Errorf("failed to unmarshal TOML: %w", err)
	}
	return nil
}

func (a *AlacrittyScheme) ToString() (string, error) {
	// Marshal the AlacrittyScheme struct into a TOML string
	data, err := toml.Marshal(a)
	if err != nil {
		return "", fmt.Errorf("failed to marshal struct to TOML: %w", err)
	}
	return string(data), nil
}
