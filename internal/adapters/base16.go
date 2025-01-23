package adapters

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Base16Scheme defines the structure for a Base16 color scheme.
type Base16Scheme struct {
	Scheme string `yaml:"scheme" abstract:"Metadata.Name"`
	Author string `yaml:"author" abstract:"Metadata.Author"`
	Base00 Color  `yaml:"base00" abstract:"SpecialColors.Background"`
	Base01 Color  `yaml:"base01" abstract:"SpecialColors.Selection"`
	Base02 Color  `yaml:"base02" abstract:"SpecialColors.Cursor"`
	Base03 Color  `yaml:"base03" abstract:"SpecialColors.CursorText"`
	Base04 Color  `yaml:"base04" abstract:"SpecialColors.SelectedText"`
	Base05 Color  `yaml:"base05" abstract:"SpecialColors.Foreground"`
	Base06 Color  `yaml:"base06"`
	Base07 Color  `yaml:"base07" abstract:"AnsiColors.White"`
	Base08 Color  `yaml:"base08" abstract:"AnsiColors.Red"`
	Base09 Color  `yaml:"base09" abstract:"AnsiColors.Yellow"`
	Base0A Color  `yaml:"base0A" abstract:"AnsiColors.Blue"`
	Base0B Color  `yaml:"base0B" abstract:"AnsiColors.Green"`
	Base0C Color  `yaml:"base0C" abstract:"AnsiColors.Cyan"`
	Base0D Color  `yaml:"base0D" abstract:"AnsiColors.BrightBlue"`
	Base0E Color  `yaml:"base0E" abstract:"AnsiColors.Magenta"`
	Base0F Color  `yaml:"base0F" abstract:"AnsiColors.BrightMagenta"`
}

// FromString converts a Base16 YAML string into an AbstractScheme.
func (rw Base16Adapter) FromString(input string) (AbstractScheme, error) {
	var base16 Base16Scheme
	err := yaml.Unmarshal([]byte(input), &base16)
	if err != nil {
		// Include a prefix with context and truncate input for readability if it's long
		return AbstractScheme{}, fmt.Errorf("failed to unmarshal Base16 YAML. Input: %.100s... Error: %w", input, err)
	}
	var abstract AbstractScheme
	PopulateAbstractScheme(base16, &abstract)

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
