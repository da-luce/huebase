package gogh

import (
	"github.com/da-luce/paletteport/internal/color"

	"gopkg.in/yaml.v3"
)

type Color = color.Color
type GoghScheme struct {
	SchemeName *string `yaml:"name" abstract:"Metadata.Name"`
	Author     *string `yaml:"author" abstract:"Metadata.Author"`
	Color01    *Color  `yaml:"color_01" abstract:"AnsiColors.Black"`
	Color02    *Color  `yaml:"color_02" abstract:"AnsiColors.Red"`
	Color03    *Color  `yaml:"color_03" abstract:"AnsiColors.Green"`
	Color04    *Color  `yaml:"color_04" abstract:"AnsiColors.Yellow"`
	Color05    *Color  `yaml:"color_05" abstract:"AnsiColors.Blue"`
	Color06    *Color  `yaml:"color_06" abstract:"AnsiColors.Magenta"`
	Color07    *Color  `yaml:"color_07" abstract:"AnsiColors.Cyan"`
	Color08    *Color  `yaml:"color_08" abstract:"AnsiColors.White"`
	Color09    *Color  `yaml:"color_09" abstract:"AnsiColors.BrightBlack"`
	Color10    *Color  `yaml:"color_10" abstract:"AnsiColors.BrightRed"`
	Color11    *Color  `yaml:"color_11" abstract:"AnsiColors.BrightGreen"`
	Color12    *Color  `yaml:"color_12" abstract:"AnsiColors.BrightYellow"`
	Color13    *Color  `yaml:"color_13" abstract:"AnsiColors.BrightBlue"`
	Color14    *Color  `yaml:"color_14" abstract:"AnsiColors.BrightMagenta"`
	Color15    *Color  `yaml:"color_15" abstract:"AnsiColors.BrightCyan"`
	Color16    *Color  `yaml:"color_16" abstract:"AnsiColors.BrightWhite"`
	Background *Color  `yaml:"background" abstract:"SpecialColors.Background"`
	Foreground *Color  `yaml:"foreground" abstract:"SpecialColors.Foreground"`
	Cursor     *Color  `yaml:"cursor" abstract:"SpecialColors.Cursor"`
}

func (rw *GoghScheme) Name() string {
	return "gogh"
}

func (rw *GoghScheme) TemplateName() string {
	return "gogh.yml.tmpl"
}

func (rw *GoghScheme) FromString(input string) error {
	err := yaml.Unmarshal([]byte(input), rw)
	if err != nil {
		return err
	}
	return nil
}
