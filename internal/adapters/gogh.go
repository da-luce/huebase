package adapters

import (
	"gopkg.in/yaml.v3"
)

type GoghScheme struct {
	Name       OptString `yaml:"name" abstract:"Metadata.Name"`
	Color01    OptColor  `yaml:"color_01" abstract:"AnsiColors.Black"`
	Color02    OptColor  `yaml:"color_02" abstract:"AnsiColors.Red"`
	Color03    OptColor  `yaml:"color_03" abstract:"AnsiColors.Green"`
	Color04    OptColor  `yaml:"color_04" abstract:"AnsiColors.Yellow"`
	Color05    OptColor  `yaml:"color_05" abstract:"AnsiColors.Blue"`
	Color06    OptColor  `yaml:"color_06" abstract:"AnsiColors.Magenta"`
	Color07    OptColor  `yaml:"color_07" abstract:"AnsiColors.Cyan"`
	Color08    OptColor  `yaml:"color_08" abstract:"AnsiColors.White"`
	Color09    OptColor  `yaml:"color_09" abstract:"AnsiColors.BrightBlack"`
	Color10    OptColor  `yaml:"color_10" abstract:"AnsiColors.BrightRed"`
	Color11    OptColor  `yaml:"color_11" abstract:"AnsiColors.BrightGreen"`
	Color12    OptColor  `yaml:"color_12" abstract:"AnsiColors.BrightYellow"`
	Color13    OptColor  `yaml:"color_13" abstract:"AnsiColors.BrightBlue"`
	Color14    OptColor  `yaml:"color_14" abstract:"AnsiColors.BrightMagenta"`
	Color15    OptColor  `yaml:"color_15" abstract:"AnsiColors.BrightCyan"`
	Color16    OptColor  `yaml:"color_16" abstract:"AnsiColors.BrightWhite"`
	Background OptColor  `yaml:"background" abstract:"SpecialColors.Background"`
	Foreground OptColor  `yaml:"foreground" abstract:"SpecialColors.Foreground"`
	Cursor     OptColor  `yaml:"cursor" abstract:"SpecialColors.Cursor"`
}

func (rw *GoghScheme) FromString(input string) error {
	err := yaml.Unmarshal([]byte(input), rw)
	if err != nil {
		return err
	}
	return nil
}

func (rw *GoghScheme) ToString() (string, error) {
	data, err := yaml.Marshal(rw)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
