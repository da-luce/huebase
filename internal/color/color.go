package color

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// Color represents a color with RGBA components.
type Color struct {
	Alpha float64
	Red   float64
	Green float64
	Blue  float64
}

// NewColor creates a new Color instance.
func NewColor(red, green, blue, alpha float64) Color {
	return Color{
		Red:   red,
		Green: green,
		Blue:  blue,
		Alpha: alpha,
	}
}

// AsString returns the color as a formatted RGBA string.
func (c Color) AsString() string {
	return fmt.Sprintf("RGBA(%.3f, %.3f, %.3f, %.3f)", c.Red, c.Green, c.Blue, c.Alpha)
}

// ToRGB converts the color to an RGB tuple.
func (c Color) ToRGB() (int, int, int) {
	return int(c.Red * 255), int(c.Green * 255), int(c.Blue * 255)
}

// ToHex converts the color to a hexadecimal string.
func (c Color) ToHex(octothorpe bool) string {
	r, g, b := c.ToRGB()
	hex := fmt.Sprintf("%02X%02X%02X", r, g, b)
	if octothorpe {
		return "#" + hex
	}
	return hex
}

// FromHex creates a Color instance from a hexadecimal string.
func FromHex(hex string) (Color, error) {

	hex = trimOctothorpe(hex)
	if len(hex) != 6 {
		return Color{}, fmt.Errorf("hex color must be 6 characters long: %s", hex)
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid red component: %w", err)
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid green component: %w", err)
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid blue component: %w", err)
	}

	return Color{
		Red:   float64(r) / 255.0,
		Green: float64(g) / 255.0,
		Blue:  float64(b) / 255.0,
		Alpha: 1.0,
	}, nil
}

// ToDict serializes the color to a dictionary-like map.
func (c Color) ToDict() map[string]string {
	return map[string]string{
		"hex": c.ToHex(true),
	}
}

// FromDict deserializes a Color from a dictionary-like map.
func FromDict(data map[string]string) (Color, error) {
	hex, exists := data["hex"]
	if !exists {
		return Color{}, fmt.Errorf("missing 'hex' key in data")
	}
	return FromHex(hex)
}

// Helper function to trim a leading '#' from a hex string.
func trimOctothorpe(hex string) string {
	if len(hex) > 0 && hex[0] == '#' {
		return hex[1:]
	}
	return hex
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

var (
	rnd  *rand.Rand
	once sync.Once
)

func getRand() *rand.Rand {
	once.Do(func() {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	return rnd
}

// RandomColor returns a random opaque Color.
func RandomColor() Color {
	r := getRand()
	return Color{
		Red:   r.Float64(),
		Green: r.Float64(),
		Blue:  r.Float64(),
		Alpha: 1.0,
	}
}
