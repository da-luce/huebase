package adapters

import (
	"bytes"

	"howett.net/plist"
)

// Temporary struct for plist parsing
type plistColor struct {
	Red   float64 `plist:"Red Component"`
	Green float64 `plist:"Green Component"`
	Blue  float64 `plist:"Blue Component"`
	Alpha float64 `plist:"Alpha Component"`
}

// iTerm scheme struct for plist serialization/deserialization
type ItermScheme struct {
	Ansi0        *Color `abstract:"AnsiColors.Black"`
	Ansi1        *Color `abstract:"AnsiColors.Red"`
	Ansi2        *Color `abstract:"AnsiColors.Green"`
	Ansi3        *Color `abstract:"AnsiColors.Yellow"`
	Ansi4        *Color `abstract:"AnsiColors.Blue"`
	Ansi5        *Color `abstract:"AnsiColors.Magenta"`
	Ansi6        *Color `abstract:"AnsiColors.Cyan"`
	Ansi7        *Color `abstract:"AnsiColors.White"`
	Ansi8        *Color `abstract:"AnsiColors.BrightBlack"`
	Ansi9        *Color `abstract:"AnsiColors.BrightRed"`
	Ansi10       *Color `abstract:"AnsiColors.BrightGreen"`
	Ansi11       *Color `abstract:"AnsiColors.BrightYellow"`
	Ansi12       *Color `abstract:"AnsiColors.BrightBlue"`
	Ansi13       *Color `abstract:"AnsiColors.BrightMagenta"`
	Ansi14       *Color `abstract:"AnsiColors.BrightCyan"`
	Ansi15       *Color `abstract:"AnsiColors.BrightWhite"`
	Background   *Color `abstract:"SpecialColors.Background"`
	Foreground   *Color `abstract:"SpecialColors.Foreground"`
	Bold         *Color `abstract:"SpecialColors.ForegroundBright"`
	Cursor       *Color `abstract:"SpecialColors.Cursor"`
	CursorText   *Color `abstract:"SpecialColors.CursorText"`
	CursorGuide  *Color `abstract:"SpecialColors.FindMatch"`
	Link         *Color `abstract:"SpecialColors.Links"`
	SelectedText *Color `abstract:"SpecialColors.SelectedText"`
	Selection    *Color `abstract:"SpecialColors.Selection"`
}

// FromString deserializes the plist data into the ItermScheme
func (rw *ItermScheme) FromString(input string) error {
	// Temporary struct to handle plist-specific parsing
	type tempItermScheme struct {
		Ansi0        *plistColor `plist:"Ansi 0 Color"`
		Ansi1        *plistColor `plist:"Ansi 1 Color"`
		Ansi2        *plistColor `plist:"Ansi 2 Color"`
		Ansi3        *plistColor `plist:"Ansi 3 Color"`
		Ansi4        *plistColor `plist:"Ansi 4 Color"`
		Ansi5        *plistColor `plist:"Ansi 5 Color"`
		Ansi6        *plistColor `plist:"Ansi 6 Color"`
		Ansi7        *plistColor `plist:"Ansi 7 Color"`
		Ansi8        *plistColor `plist:"Ansi 8 Color"`
		Ansi9        *plistColor `plist:"Ansi 9 Color"`
		Ansi10       *plistColor `plist:"Ansi 10 Color"`
		Ansi11       *plistColor `plist:"Ansi 11 Color"`
		Ansi12       *plistColor `plist:"Ansi 12 Color"`
		Ansi13       *plistColor `plist:"Ansi 13 Color"`
		Ansi14       *plistColor `plist:"Ansi 14 Color"`
		Ansi15       *plistColor `plist:"Ansi 15 Color"`
		Background   *plistColor `plist:"Background Color"`
		Foreground   *plistColor `plist:"Foreground Color"`
		Bold         *plistColor `plist:"Bold Color"`
		Cursor       *plistColor `plist:"Cursor Color"`
		CursorText   *plistColor `plist:"Cursor Text Color"`
		CursorGuide  *plistColor `plist:"Cursor Guide Color"`
		Link         *plistColor `plist:"Link Color"`
		SelectedText *plistColor `plist:"Selected Text Color"`
		Selection    *plistColor `plist:"Selection Color"`
	}

	temp := &tempItermScheme{}
	decoder := plist.NewDecoder(bytes.NewReader([]byte(input)))
	if err := decoder.Decode(temp); err != nil {
		return err
	}

	// Convert plistColor to Color
	rw.Ansi0 = convertToColor(temp.Ansi0)
	rw.Ansi1 = convertToColor(temp.Ansi1)
	rw.Ansi2 = convertToColor(temp.Ansi2)
	rw.Ansi3 = convertToColor(temp.Ansi3)
	rw.Ansi4 = convertToColor(temp.Ansi4)
	rw.Ansi5 = convertToColor(temp.Ansi5)
	rw.Ansi6 = convertToColor(temp.Ansi6)
	rw.Ansi7 = convertToColor(temp.Ansi7)
	rw.Ansi8 = convertToColor(temp.Ansi8)
	rw.Ansi9 = convertToColor(temp.Ansi9)
	rw.Ansi10 = convertToColor(temp.Ansi10)
	rw.Ansi11 = convertToColor(temp.Ansi11)
	rw.Ansi12 = convertToColor(temp.Ansi12)
	rw.Ansi13 = convertToColor(temp.Ansi13)
	rw.Ansi14 = convertToColor(temp.Ansi14)
	rw.Ansi15 = convertToColor(temp.Ansi15)
	rw.Background = convertToColor(temp.Background)
	rw.Foreground = convertToColor(temp.Foreground)
	rw.Bold = convertToColor(temp.Bold)
	rw.Cursor = convertToColor(temp.Cursor)
	rw.CursorText = convertToColor(temp.CursorText)
	rw.CursorGuide = convertToColor(temp.CursorGuide)
	rw.Link = convertToColor(temp.Link)
	rw.SelectedText = convertToColor(temp.SelectedText)
	rw.Selection = convertToColor(temp.Selection)

	return nil
}

// ToString serializes the ItermScheme into a plist string
func (rw *ItermScheme) ToString() (string, error) {
	// Temporary struct to handle plist-specific serialization
	type tempItermScheme struct {
		Ansi0        *plistColor `plist:"Ansi 0 Color"`
		Ansi1        *plistColor `plist:"Ansi 1 Color"`
		Ansi2        *plistColor `plist:"Ansi 2 Color"`
		Ansi3        *plistColor `plist:"Ansi 3 Color"`
		Ansi4        *plistColor `plist:"Ansi 4 Color"`
		Ansi5        *plistColor `plist:"Ansi 5 Color"`
		Ansi6        *plistColor `plist:"Ansi 6 Color"`
		Ansi7        *plistColor `plist:"Ansi 7 Color"`
		Ansi8        *plistColor `plist:"Ansi 8 Color"`
		Ansi9        *plistColor `plist:"Ansi 9 Color"`
		Ansi10       *plistColor `plist:"Ansi 10 Color"`
		Ansi11       *plistColor `plist:"Ansi 11 Color"`
		Ansi12       *plistColor `plist:"Ansi 12 Color"`
		Ansi13       *plistColor `plist:"Ansi 13 Color"`
		Ansi14       *plistColor `plist:"Ansi 14 Color"`
		Ansi15       *plistColor `plist:"Ansi 15 Color"`
		Background   *plistColor `plist:"Background Color"`
		Foreground   *plistColor `plist:"Foreground Color"`
		Bold         *plistColor `plist:"Bold Color"`
		Cursor       *plistColor `plist:"Cursor Color"`
		CursorText   *plistColor `plist:"Cursor Text Color"`
		CursorGuide  *plistColor `plist:"Cursor Guide Color"`
		Link         *plistColor `plist:"Link Color"`
		SelectedText *plistColor `plist:"Selected Text Color"`
		Selection    *plistColor `plist:"Selection Color"`
	}

	temp := &tempItermScheme{
		Ansi0:        convertToPlistColor(rw.Ansi0),
		Ansi1:        convertToPlistColor(rw.Ansi1),
		Ansi2:        convertToPlistColor(rw.Ansi2),
		Ansi3:        convertToPlistColor(rw.Ansi3),
		Ansi4:        convertToPlistColor(rw.Ansi4),
		Ansi5:        convertToPlistColor(rw.Ansi5),
		Ansi6:        convertToPlistColor(rw.Ansi6),
		Ansi7:        convertToPlistColor(rw.Ansi7),
		Ansi8:        convertToPlistColor(rw.Ansi8),
		Ansi9:        convertToPlistColor(rw.Ansi9),
		Ansi10:       convertToPlistColor(rw.Ansi10),
		Ansi11:       convertToPlistColor(rw.Ansi11),
		Ansi12:       convertToPlistColor(rw.Ansi12),
		Ansi13:       convertToPlistColor(rw.Ansi13),
		Ansi14:       convertToPlistColor(rw.Ansi14),
		Ansi15:       convertToPlistColor(rw.Ansi15),
		Background:   convertToPlistColor(rw.Background),
		Foreground:   convertToPlistColor(rw.Foreground),
		Bold:         convertToPlistColor(rw.Bold),
		Cursor:       convertToPlistColor(rw.Cursor),
		CursorText:   convertToPlistColor(rw.CursorText),
		CursorGuide:  convertToPlistColor(rw.CursorGuide),
		Link:         convertToPlistColor(rw.Link),
		SelectedText: convertToPlistColor(rw.SelectedText),
		Selection:    convertToPlistColor(rw.Selection),
	}

	var buf bytes.Buffer
	encoder := plist.NewEncoder(&buf)
	err := encoder.Encode(temp)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Helper functions for conversion
func convertToColor(pc *plistColor) *Color {
	if pc == nil {
		return nil
	}
	return &Color{
		Red:   pc.Red,
		Green: pc.Green,
		Blue:  pc.Blue,
		Alpha: pc.Alpha,
	}
}

func convertToPlistColor(c *Color) *plistColor {
	if c == nil {
		return nil
	}
	return &plistColor{
		Red:   c.Red,
		Green: c.Green,
		Blue:  c.Blue,
		Alpha: c.Alpha,
	}
}
