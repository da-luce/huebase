package adapters

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

// Primary colors
type Primary struct {
	Background OptColor `toml:"background" abstract:"SpecialColors.Background"`
	Foreground OptColor `toml:"foreground" abstract:"SpecialColors.Foreground"`
}

// Cursor colors
type Cursor struct {
	Cursor OptColor `toml:"cursor" abstract:"SpecialColors.Cursor"`
	Text   OptColor `toml:"text" abstract:"SpecialColors.CursorText"`
}

// Normal ANSI colors
type Normal struct {
	Black   OptColor `toml:"black" abstract:"AnsiColors.Black"`
	Blue    OptColor `toml:"blue" abstract:"AnsiColors.Blue"`
	Cyan    OptColor `toml:"cyan" abstract:"AnsiColors.Cyan"`
	Green   OptColor `toml:"green" abstract:"AnsiColors.Green"`
	Magenta OptColor `toml:"magenta" abstract:"AnsiColors.Magenta"`
	Red     OptColor `toml:"red" abstract:"AnsiColors.Red"`
	White   OptColor `toml:"white" abstract:"AnsiColors.White"`
	Yellow  OptColor `toml:"yellow" abstract:"AnsiColors.Yellow"`
}

// Bright ANSI colors
type Bright struct {
	Black   OptColor `toml:"black" abstract:"AnsiColors.BrightBlack"`
	Blue    OptColor `toml:"blue" abstract:"AnsiColors.BrightBlue"`
	Cyan    OptColor `toml:"cyan" abstract:"AnsiColors.BrightCyan"`
	Green   OptColor `toml:"green" abstract:"AnsiColors.BrightGreen"`
	Magenta OptColor `toml:"magenta" abstract:"AnsiColors.BrightMagenta"`
	Red     OptColor `toml:"red" abstract:"AnsiColors.BrightRed"`
	White   OptColor `toml:"white" abstract:"AnsiColors.BrightWhite"`
	Yellow  OptColor `toml:"yellow" abstract:"AnsiColors.BrightYellow"`
}

// Selection colors
type Selection struct {
	Background OptColor `toml:"background" abstract:"SpecialColors.Selection"`
	Text       OptColor `toml:"text" abstract:"SpecialColors.SelectedText"`
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
	err := toml.Unmarshal([]byte(input), a)
	if err != nil {

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
