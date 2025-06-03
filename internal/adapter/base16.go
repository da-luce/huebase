package adapter

import (
	"gopkg.in/yaml.v3"
)

type Base16Scheme struct {
	Scheme *string `yaml:"scheme"`
	Author *string `yaml:"author"`
	Base00 *Color  `yaml:"base00"`
	Base01 *Color  `yaml:"base01"`
	Base02 *Color  `yaml:"base02"`
	Base03 *Color  `yaml:"base03"`
	Base04 *Color  `yaml:"base04"`
	Base05 *Color  `yaml:"base05"`
	Base06 *Color  `yaml:"base06"`
	Base07 *Color  `yaml:"base07"`
	Base08 *Color  `yaml:"base08"`
	Base09 *Color  `yaml:"base09"`
	Base0A *Color  `yaml:"base0A"`
	Base0B *Color  `yaml:"base0B"`
	Base0C *Color  `yaml:"base0C"`
	Base0D *Color  `yaml:"base0D"`
	Base0E *Color  `yaml:"base0E"`
	Base0F *Color  `yaml:"base0F"`
}

func (rw *Base16Scheme) FromString(input string) error {
	err := yaml.Unmarshal([]byte(input), rw)
	if err != nil {
		return err
	}
	return nil
}

func (rw *Base16Scheme) MappingPath() string {
	return "./base16.yaml"
}

func (rw *Base16Scheme) TemplatePath() string {
	return "../templates/base16.yaml"
}

func (rw *Base16Scheme) Name() string {
	return "base16"
}
