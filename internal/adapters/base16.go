package adapters

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Base16Scheme defines the structure for a Base16 color scheme.
type Base16Scheme struct {
	Scheme string `yaml:"scheme"`
	Author string `yaml:"author"`
	Base00 Color  `yaml:"base00"`
	Base01 Color  `yaml:"base01"`
	Base02 Color  `yaml:"base02"`
	Base03 Color  `yaml:"base03"`
	Base04 Color  `yaml:"base04"`
	Base05 Color  `yaml:"base05"`
	Base06 Color  `yaml:"base06"`
	Base07 Color  `yaml:"base07"`
	Base08 Color  `yaml:"base08"`
	Base09 Color  `yaml:"base09"`
	Base0A Color  `yaml:"base0A"`
	Base0B Color  `yaml:"base0B"`
	Base0C Color  `yaml:"base0C"`
	Base0D Color  `yaml:"base0D"`
	Base0E Color  `yaml:"base0E"`
	Base0F Color  `yaml:"base0F"`
}

// FromString converts a Base16 YAML string into an AbstractScheme.
func (rw Base16Adapter) FromString(input string) (AbstractScheme, error) {
	var base16 Base16Scheme
	err := yaml.Unmarshal([]byte(input), &base16)
	if err != nil {
		// Include a prefix with context and truncate input for readability if it's long
		return AbstractScheme{}, fmt.Errorf("failed to unmarshal Base16 YAML. Input: %.100s... Error: %w", input, err)
	}

	abstract := AbstractScheme{
		AnsiColors: AnsiColors{
			Black:         base16.Base00,
			Red:           base16.Base08,
			Green:         base16.Base0B,
			Yellow:        base16.Base0A,
			Blue:          base16.Base0D,
			Magenta:       base16.Base0E,
			Cyan:          base16.Base0C,
			White:         base16.Base05,
			BrightBlack:   base16.Base03,
			BrightRed:     base16.Base09,
			BrightGreen:   base16.Base0B,
			BrightYellow:  base16.Base0A,
			BrightBlue:    base16.Base0D,
			BrightMagenta: base16.Base0E,
			BrightCyan:    base16.Base0C,
			BrightWhite:   base16.Base07,
		},
		SpecialColors: SpecialColors{
			Foreground:   base16.Base05,
			Background:   base16.Base00,
			Cursor:       base16.Base05,
			CursorText:   base16.Base00,
			Selection:    base16.Base02,
			SelectedText: base16.Base05,
			Links:        base16.Base0D,
			FindMatch:    base16.Base09,
		},
		Metadata: Meta{
			Name:   base16.Scheme,
			Author: base16.Author,
		},
	}
	return abstract, nil
}

// ToString converts an AbstractScheme into a Base16 YAML string.
func (rw Base16Adapter) ToString(theme AbstractScheme) (string, error) {
	base16 := Base16Scheme{
		Scheme: theme.Metadata.Name,
		Author: theme.Metadata.Author,
		Base00: theme.SpecialColors.Background,
		Base01: theme.SpecialColors.Selection,
		Base02: theme.SpecialColors.Selection,
		Base03: theme.AnsiColors.BrightBlack,
		Base04: theme.SpecialColors.Cursor,
		Base05: theme.SpecialColors.Foreground,
		Base06: theme.AnsiColors.White,
		Base07: theme.AnsiColors.BrightWhite,
		Base08: theme.AnsiColors.Red,
		Base09: theme.AnsiColors.BrightRed,
		Base0A: theme.AnsiColors.Yellow,
		Base0B: theme.AnsiColors.Green,
		Base0C: theme.AnsiColors.Cyan,
		Base0D: theme.AnsiColors.Blue,
		Base0E: theme.AnsiColors.Magenta,
		Base0F: theme.AnsiColors.BrightMagenta,
	}

	data, err := yaml.Marshal(base16)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
