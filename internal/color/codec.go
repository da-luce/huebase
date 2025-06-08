package color

import (
	"encoding/json"
	"fmt"
	"html/template"
)

// Implements the plist.Marshaler and plist.Unmarshaler behavior via plist tags
func (c *Color) UnmarshalPlist(unmarshal func(interface{}) error) error {
	var dict map[string]float64
	if err := unmarshal(&dict); err != nil {
		return err
	}

	// Extract color components safely
	var ok bool
	if c.Red, ok = dict["Red Component"]; !ok {
		return fmt.Errorf("missing Red Component")
	}
	if c.Green, ok = dict["Green Component"]; !ok {
		return fmt.Errorf("missing Green Component")
	}
	if c.Blue, ok = dict["Blue Component"]; !ok {
		return fmt.Errorf("missing Blue Component")
	}
	if alpha, found := dict["Alpha Component"]; found {
		c.Alpha = alpha
	} else {
		c.Alpha = 1.0 // default alpha if not provided
	}

	return nil
}

func (c *Color) ToITermXML() template.HTML {
	var r, g, b float64
	if c == nil {
		r, g, b = 0, 0, 0 // or whatever default you want
	} else {
		r, g, b = c.Red, c.Green, c.Blue
	}

	dict := fmt.Sprintf(
		`<dict>
    <key>Blue Component</key>
    <real>%f</real>
    <key>Green Component</key>
    <real>%f</real>
    <key>Red Component</key>
    <real>%f</real>
</dict>`, b, g, r)
	return template.HTML(dict)
}

// UnmarshalYAML allows YAML to deserialize directly into the Color type.
func (c *Color) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var colorStr string

	if err := unmarshal(&colorStr); err != nil {
		return err
	}
	color, err := FromHex(colorStr)
	if err != nil {
		return err
	}
	*c = color
	return nil
}

// UnmarshalJSON allows JSON to deserialize directly into the Color type from a hex string.
func (c *Color) UnmarshalJSON(data []byte) error {
	// JSON string comes with quotes, so unquote it first
	var colorStr string
	if err := json.Unmarshal(data, &colorStr); err != nil {
		return err
	}

	color, err := FromHex(colorStr)
	if err != nil {
		return fmt.Errorf("invalid color hex: %w", err)
	}

	*c = color
	return nil
}

// MarshalYAML allows YAML to serialize the Color type.
func (c Color) MarshalYAML() (interface{}, error) {
	return c.ToHex(true), nil
}

// UnmarshalTOML allows TOML to deserialize directly into the Color type.
func (c *Color) UnmarshalTOML(data interface{}) error {

	colorStr, ok := data.(string)

	if !ok {
		return fmt.Errorf("invalid color format: expected a string, got %T (value: %v)", data, data)
	}
	color, err := FromHex(colorStr)
	if err != nil {
		return err
	}
	*c = color
	return nil
}

func (c *Color) UnmarshalText(text []byte) error {

	color, err := FromHex(string(text))
	if err != nil {
		return err
	}
	*c = color
	return nil
}

func (c *Color) MarshalText() (text []byte, err error) {
	return []byte(c.ToHex(true)), nil
}

// MarshalTOML allows TOML to serialize the Color type.
func (c Color) MarshalTOML() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", c.ToHex(true))), nil
}
