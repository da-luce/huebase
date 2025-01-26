package adapters

import (
	"gopkg.in/yaml.v3"
)

type Base16Scheme struct {
	Scheme OptString `yaml:"scheme" abstract:"Metadata.Name"`
	Author OptString `yaml:"author" abstract:"Metadata.Author"`
	Base00 OptColor  `yaml:"base00" abstract:"SpecialColors.Background"`
	Base01 OptColor  `yaml:"base01" abstract:"SpecialColors.Selection"`
	Base02 OptColor  `yaml:"base02" abstract:"SpecialColors.Cursor"`
	Base03 OptColor  `yaml:"base03" abstract:"SpecialColors.CursorText"`
	Base04 OptColor  `yaml:"base04" abstract:"SpecialColors.SelectedText"`
	Base05 OptColor  `yaml:"base05" abstract:"SpecialColors.Foreground"`
	Base06 OptColor  `yaml:"base06" abstract:"SpecialColors.ForegroundBright"`
	Base07 OptColor  `yaml:"base07" abstract:"AnsiColors.White"`
	Base08 OptColor  `yaml:"base08" abstract:"AnsiColors.Red"`
	Base09 OptColor  `yaml:"base09" abstract:"AnsiColors.Yellow"`
	Base0A OptColor  `yaml:"base0A" abstract:"AnsiColors.Blue"`
	Base0B OptColor  `yaml:"base0B" abstract:"AnsiColors.Green"`
	Base0C OptColor  `yaml:"base0C" abstract:"AnsiColors.Cyan"`
	Base0D OptColor  `yaml:"base0D" abstract:"AnsiColors.BrightBlue"`
	Base0E OptColor  `yaml:"base0E" abstract:"AnsiColors.Magenta"`
	Base0F OptColor  `yaml:"base0F" abstract:"AnsiColors.BrightMagenta"`
}

func (rw *Base16Scheme) FromString(input string) error {
	err := yaml.Unmarshal([]byte(input), rw)
	if err != nil {
		return err
	}
	return nil
}

func (rw *Base16Scheme) ToString() (string, error) {
	// Marshal the Base16Scheme struct to a YAML string
	data, err := yaml.Marshal(rw)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (rw *Base16Scheme) ToAbstract() (AbstractScheme, error) {
	return ToAbstractDefault(rw)
}

func (rw *Base16Scheme) FromAbstract(abstract *AbstractScheme) error {
	return FromAbstractDefault(abstract, rw)
}
