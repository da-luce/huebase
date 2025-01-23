package adapters

import "github.com/BurntSushi/toml"

// Use a TOML library for parsing and encoding

// AlacrittyScheme defines the structure for an Alacritty color scheme.
type AlacrittyScheme struct {
	Primary struct {
		Background Color `toml:"background"`
		Foreground Color `toml:"foreground"`
	} `toml:"colors.primary"`

	Cursor struct {
		Cursor Color `toml:"cursor"`
		Text   Color `toml:"text"`
	} `toml:"colors.cursor"`

	Normal struct {
		Black   Color `toml:"black"`
		Blue    Color `toml:"blue"`
		Cyan    Color `toml:"cyan"`
		Green   Color `toml:"green"`
		Magenta Color `toml:"magenta"`
		Red     Color `toml:"red"`
		White   Color `toml:"white"`
		Yellow  Color `toml:"yellow"`
	} `toml:"colors.normal"`

	Bright struct {
		Black   Color `toml:"black"`
		Blue    Color `toml:"blue"`
		Cyan    Color `toml:"cyan"`
		Green   Color `toml:"green"`
		Magenta Color `toml:"magenta"`
		Red     Color `toml:"red"`
		White   Color `toml:"white"`
		Yellow  Color `toml:"yellow"`
	} `toml:"colors.bright"`

	Selection struct {
		Background Color `toml:"background"`
		Text       Color `toml:"text"`
	} `toml:"colors.selection"`
}

// FromString parses an Alacritty TOML string into an AbstractScheme.
func (rw AlacrittyAdapter) FromString(input string) (AbstractScheme, error) {
	var alacritty AlacrittyScheme
	err := toml.Unmarshal([]byte(input), &alacritty)
	if err != nil {
		return AbstractScheme{}, err
	}

	abstract := AbstractScheme{
		AnsiColors: AnsiColors{
			Black:         alacritty.Normal.Black,
			Red:           alacritty.Normal.Red,
			Green:         alacritty.Normal.Green,
			Yellow:        alacritty.Normal.Yellow,
			Blue:          alacritty.Normal.Blue,
			Magenta:       alacritty.Normal.Magenta,
			Cyan:          alacritty.Normal.Cyan,
			White:         alacritty.Normal.White,
			BrightBlack:   alacritty.Bright.Black,
			BrightRed:     alacritty.Bright.Red,
			BrightGreen:   alacritty.Bright.Green,
			BrightYellow:  alacritty.Bright.Yellow,
			BrightBlue:    alacritty.Bright.Blue,
			BrightMagenta: alacritty.Bright.Magenta,
			BrightCyan:    alacritty.Bright.Cyan,
			BrightWhite:   alacritty.Bright.White,
		},
		SpecialColors: SpecialColors{
			Foreground:   alacritty.Primary.Foreground,
			Background:   alacritty.Primary.Background,
			Cursor:       alacritty.Cursor.Cursor,
			CursorText:   alacritty.Cursor.Text,
			Selection:    alacritty.Selection.Background,
			SelectedText: alacritty.Selection.Text,
		},
		Metadata: Meta{
			Name:   "Alacritty Scheme", // Metadata can be set manually if not part of the TOML
			Author: "Unknown",          // Optional: Extract from comments if needed
		},
	}
	return abstract, nil
}

// ToString serializes an AbstractScheme into an Alacritty TOML string.
func (rw AlacrittyAdapter) ToString(theme AbstractScheme) (string, error) {
	alacritty := AlacrittyScheme{
		Primary: struct {
			Background Color `toml:"background"`
			Foreground Color `toml:"foreground"`
		}{
			Background: theme.SpecialColors.Background,
			Foreground: theme.SpecialColors.Foreground,
		},
		Cursor: struct {
			Cursor Color `toml:"cursor"`
			Text   Color `toml:"text"`
		}{
			Cursor: theme.SpecialColors.Cursor,
			Text:   theme.SpecialColors.CursorText,
		},
		Normal: struct {
			Black   Color `toml:"black"`
			Blue    Color `toml:"blue"`
			Cyan    Color `toml:"cyan"`
			Green   Color `toml:"green"`
			Magenta Color `toml:"magenta"`
			Red     Color `toml:"red"`
			White   Color `toml:"white"`
			Yellow  Color `toml:"yellow"`
		}{
			Black:   theme.AnsiColors.Black,
			Red:     theme.AnsiColors.Red,
			Green:   theme.AnsiColors.Green,
			Yellow:  theme.AnsiColors.Yellow,
			Blue:    theme.AnsiColors.Blue,
			Magenta: theme.AnsiColors.Magenta,
			Cyan:    theme.AnsiColors.Cyan,
			White:   theme.AnsiColors.White,
		},
		Bright: struct {
			Black   Color `toml:"black"`
			Blue    Color `toml:"blue"`
			Cyan    Color `toml:"cyan"`
			Green   Color `toml:"green"`
			Magenta Color `toml:"magenta"`
			Red     Color `toml:"red"`
			White   Color `toml:"white"`
			Yellow  Color `toml:"yellow"`
		}{
			Black:   theme.AnsiColors.BrightBlack,
			Red:     theme.AnsiColors.BrightRed,
			Green:   theme.AnsiColors.BrightGreen,
			Yellow:  theme.AnsiColors.BrightYellow,
			Blue:    theme.AnsiColors.BrightBlue,
			Magenta: theme.AnsiColors.BrightMagenta,
			Cyan:    theme.AnsiColors.BrightCyan,
			White:   theme.AnsiColors.BrightWhite,
		},
		Selection: struct {
			Background Color `toml:"background"`
			Text       Color `toml:"text"`
		}{
			Background: theme.SpecialColors.Selection,
			Text:       theme.SpecialColors.SelectedText,
		},
	}

	data, err := toml.Marshal(alacritty)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
