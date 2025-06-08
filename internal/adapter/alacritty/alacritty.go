package alacritty

import (
	"github.com/da-luce/paletteport/internal/color"

	"github.com/pelletier/go-toml/v2"
)

type Color = color.Color
type Primary struct {
	Background *Color `toml:"background" abstract:"SpecialColors.Background"`
	Foreground *Color `toml:"foreground" abstract:"SpecialColors.Foreground"`
}

type Cursor struct {
	Cursor *Color `toml:"cursor" abstract:"SpecialColors.Cursor"`
	Text   *Color `toml:"text" abstract:"SpecialColors.CursorText"`
}

type Normal struct {
	Black   *Color `toml:"black" abstract:"AnsiColors.Black"`
	Blue    *Color `toml:"blue" abstract:"AnsiColors.Blue"`
	Cyan    *Color `toml:"cyan" abstract:"AnsiColors.Cyan"`
	Green   *Color `toml:"green" abstract:"AnsiColors.Green"`
	Magenta *Color `toml:"magenta" abstract:"AnsiColors.Magenta"`
	Red     *Color `toml:"red" abstract:"AnsiColors.Red"`
	White   *Color `toml:"white" abstract:"AnsiColors.White"`
	Yellow  *Color `toml:"yellow" abstract:"AnsiColors.Yellow"`
}

type Bright struct {
	Black   *Color `toml:"black" abstract:"AnsiColors.BrightBlack"`
	Blue    *Color `toml:"blue" abstract:"AnsiColors.BrightBlue"`
	Cyan    *Color `toml:"cyan" abstract:"AnsiColors.BrightCyan"`
	Green   *Color `toml:"green" abstract:"AnsiColors.BrightGreen"`
	Magenta *Color `toml:"magenta" abstract:"AnsiColors.BrightMagenta"`
	Red     *Color `toml:"red" abstract:"AnsiColors.BrightRed"`
	White   *Color `toml:"white" abstract:"AnsiColors.BrightWhite"`
	Yellow  *Color `toml:"yellow" abstract:"AnsiColors.BrightYellow"`
}

type Selection struct {
	Background *Color `toml:"background" abstract:"SpecialColors.Selection"`
	Text       *Color `toml:"text" abstract:"SpecialColors.SelectedText"`
}

type Colors struct {
	Primary   Primary   `toml:"primary"`
	Cursor    Cursor    `toml:"cursor"`
	Normal    Normal    `toml:"normal"`
	Bright    Bright    `toml:"bright"`
	Selection Selection `toml:"selection"`
}

type AlacrittyScheme struct {
	Colors Colors `toml:"colors"`
}

func (rw *AlacrittyScheme) Name() string {
	return "alacritty"
}

func (rw *AlacrittyScheme) TemplateName() string {
	return "alacritty.toml.tmpl"
}

func (a *AlacrittyScheme) FromString(input string) error {
	err := toml.Unmarshal([]byte(input), a)
	if err != nil {
		return err
	}
	return nil
}
