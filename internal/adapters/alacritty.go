package adapters

import "github.com/pelletier/go-toml"

// AlacrittyScheme defines the structure for an Alacritty color scheme.
type Primary struct {
	Background Color `toml:"background"`
	Foreground Color `toml:"foreground"`
}

type Cursor struct {
	Cursor Color `toml:"cursor"`
	Text   Color `toml:"text"`
}

type Normal struct {
	Black   Color `toml:"black"`
	Blue    Color `toml:"blue"`
	Cyan    Color `toml:"cyan"`
	Green   Color `toml:"green"`
	Magenta Color `toml:"magenta"`
	Red     Color `toml:"red"`
	White   Color `toml:"white"`
	Yellow  Color `toml:"yellow"`
}

type Bright struct {
	Black   Color `toml:"black"`
	Blue    Color `toml:"blue"`
	Cyan    Color `toml:"cyan"`
	Green   Color `toml:"green"`
	Magenta Color `toml:"magenta"`
	Red     Color `toml:"red"`
	White   Color `toml:"white"`
	Yellow  Color `toml:"yellow"`
}

type Selection struct {
	Background Color `toml:"background"`
	Text       Color `toml:"text"`
}

// Combine into a Colors struct
type Colors struct {
	Primary   Primary   `toml:"primary"`
	Cursor    Cursor    `toml:"cursor"`
	Normal    Normal    `toml:"normal"`
	Bright    Bright    `toml:"colors.bright"`
	Selection Selection `toml:"selection"`
}

// Final AlacrittyScheme struct
type AlacrittyScheme struct {
	Colors Colors `toml:"colors"`
}

// FromString parses an Alacritty TOML string into an AbstractScheme.
func (rw AlacrittyAdapter) FromString(input string) (AbstractScheme, error) {
	var alacritty AlacrittyScheme
	if err := toml.Unmarshal([]byte(input), &alacritty); err != nil {
		return AbstractScheme{}, err
	}

	// Map parsed AlacrittyScheme to AbstractScheme
	abstract := AbstractScheme{
		AnsiColors: AnsiColors{
			Black:         alacritty.Colors.Normal.Black,
			Red:           alacritty.Colors.Normal.Red,
			Green:         alacritty.Colors.Normal.Green,
			Yellow:        alacritty.Colors.Normal.Yellow,
			Blue:          alacritty.Colors.Normal.Blue,
			Magenta:       alacritty.Colors.Normal.Magenta,
			Cyan:          alacritty.Colors.Normal.Cyan,
			White:         alacritty.Colors.Normal.White,
			BrightBlack:   alacritty.Colors.Bright.Black,
			BrightRed:     alacritty.Colors.Bright.Red,
			BrightGreen:   alacritty.Colors.Bright.Green,
			BrightYellow:  alacritty.Colors.Bright.Yellow,
			BrightBlue:    alacritty.Colors.Bright.Blue,
			BrightMagenta: alacritty.Colors.Bright.Magenta,
			BrightCyan:    alacritty.Colors.Bright.Cyan,
			BrightWhite:   alacritty.Colors.Bright.White,
		},
		SpecialColors: SpecialColors{
			Foreground:   alacritty.Colors.Primary.Foreground,
			Background:   alacritty.Colors.Primary.Background,
			Cursor:       alacritty.Colors.Cursor.Cursor,
			CursorText:   alacritty.Colors.Cursor.Text,
			Selection:    alacritty.Colors.Selection.Background,
			SelectedText: alacritty.Colors.Selection.Text,
		},
		Metadata: Meta{
			Name:   "Unknown", // Metadata can be set manually if not part of the TOML
			Author: "Unknown", // Optional: Extract from comments if needed
		},
	}
	return abstract, nil
}

// ToString serializes an AbstractScheme into an Alacritty TOML string.
func (rw AlacrittyAdapter) ToString(theme AbstractScheme) (string, error) {
	// Map AbstractScheme to AlacrittyScheme
	alacritty := AlacrittyScheme{
		Colors{
			Primary{
				Background: theme.SpecialColors.Background,
				Foreground: theme.SpecialColors.Foreground,
			},
			Cursor{
				Cursor: theme.SpecialColors.Cursor,
				Text:   theme.SpecialColors.CursorText,
			},
			Normal{
				Black:   theme.AnsiColors.Black,
				Red:     theme.AnsiColors.Red,
				Green:   theme.AnsiColors.Green,
				Yellow:  theme.AnsiColors.Yellow,
				Blue:    theme.AnsiColors.Blue,
				Magenta: theme.AnsiColors.Magenta,
				Cyan:    theme.AnsiColors.Cyan,
				White:   theme.AnsiColors.White,
			},
			Bright{
				Black:   theme.AnsiColors.BrightBlack,
				Red:     theme.AnsiColors.BrightRed,
				Green:   theme.AnsiColors.BrightGreen,
				Yellow:  theme.AnsiColors.BrightYellow,
				Blue:    theme.AnsiColors.BrightBlue,
				Magenta: theme.AnsiColors.BrightMagenta,
				Cyan:    theme.AnsiColors.BrightCyan,
				White:   theme.AnsiColors.BrightWhite,
			},
			Selection{
				Background: theme.SpecialColors.Selection,
				Text:       theme.SpecialColors.SelectedText,
			},
		},
	}

	// Serialize to TOML
	data, err := toml.Marshal(alacritty)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
