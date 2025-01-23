package adapters

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

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

func (rw Base16Scheme) FromString(input string) (Adapter, error) {

	var base16 Base16Scheme
	err := yaml.Unmarshal([]byte(input), &base16)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Base16 YAML: %w", err)
	}
	return base16, nil
}

func (adapter Base16Scheme) ToString(input interface{}) (string, error) {

	data, err := yaml.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal struct to YAML: %w", err)
	}
	return string(data), nil
}
