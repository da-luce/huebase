package adapters

import (
	"encoding/xml"
	"fmt"
)

// iTermScheme represents the structure of an iTerm2 color scheme.
type iTermScheme struct {
	Name OptString `xml:"-" abstract:"Metadata.Name"` // Metadata.Name for the scheme name
	// AnsiColors        [16]Color `xml:"dict>key>Ansi" abstract:"AnsiColors"` // Maps to AbstractScheme.AnsiColors
	BackgroundColor   Color `xml:"dict>key>Background Color" abstract:"SpecialColors.Background"`
	ForegroundColor   Color `xml:"dict>key>Foreground Color" abstract:"SpecialColors.Foreground"`
	CursorColor       Color `xml:"dict>key>Cursor Color" abstract:"SpecialColors.Cursor"`
	SelectionColor    Color `xml:"dict>key>Selection Color" abstract:"SpecialColors.Selection"`
	SelectedTextColor Color `xml:"dict>key>Selected Text Color" abstract:"SpecialColors.SelectedText"`
	CursorTextColor   Color `xml:"dict>key>Cursor Text Color" abstract:"SpecialColors.CursorText"`
}

// FromString unmarshals the iTerm2 color scheme from a plist XML string.
func (rw *iTermScheme) FromString(input string) error {
	err := xml.Unmarshal([]byte(input), rw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal iTerm2 plist: %w", err)
	}
	return nil
}

// ToString marshals the iTerm2 color scheme into a plist XML string.
func (rw *iTermScheme) ToString() (string, error) {
	data, err := xml.MarshalIndent(rw, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal iTerm2 plist: %w", err)
	}
	return xml.Header + string(data), nil
}
