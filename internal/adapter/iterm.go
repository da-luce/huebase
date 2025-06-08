package adapter

import (
	"bytes"

	"howett.net/plist"
)

// iTerm scheme struct for plist serialization/deserialization
type ItermScheme struct {
	Ansi0        *Color `plist:"Ansi 0 Color" abstract:"AnsiColors.Black"`
	Ansi1        *Color `plist:"Ansi 1 Color" abstract:"AnsiColors.Red"`
	Ansi2        *Color `plist:"Ansi 2 Color" abstract:"AnsiColors.Green"`
	Ansi3        *Color `plist:"Ansi 3 Color" abstract:"AnsiColors.Yellow"`
	Ansi4        *Color `plist:"Ansi 4 Color" abstract:"AnsiColors.Blue"`
	Ansi5        *Color `plist:"Ansi 5 Color" abstract:"AnsiColors.Magenta"`
	Ansi6        *Color `plist:"Ansi 6 Color" abstract:"AnsiColors.Cyan"`
	Ansi7        *Color `plist:"Ansi 7 Color" abstract:"AnsiColors.White"`
	Ansi8        *Color `plist:"Ansi 8 Color" abstract:"AnsiColors.BrightBlack"`
	Ansi9        *Color `plist:"Ansi 9 Color" abstract:"AnsiColors.BrightRed"`
	Ansi10       *Color `plist:"Ansi 10 Color" abstract:"AnsiColors.BrightGreen"`
	Ansi11       *Color `plist:"Ansi 11 Color" abstract:"AnsiColors.BrightYellow"`
	Ansi12       *Color `plist:"Ansi 12 Color" abstract:"AnsiColors.BrightBlue"`
	Ansi13       *Color `plist:"Ansi 13 Color" abstract:"AnsiColors.BrightMagenta"`
	Ansi14       *Color `plist:"Ansi 14 Color" abstract:"AnsiColors.BrightCyan"`
	Ansi15       *Color `plist:"Ansi 15 Color" abstract:"AnsiColors.BrightWhite"`
	Background   *Color `plist:"Background Color" abstract:"SpecialColors.Background"`
	Foreground   *Color `plist:"Foreground Color" abstract:"SpecialColors.Foreground"`
	Bold         *Color `plist:"Bold Color" abstract:"SpecialColors.ForegroundBright"`
	Cursor       *Color `plist:"Cursor Color" abstract:"SpecialColors.Cursor"`
	CursorText   *Color `plist:"Cursor Text Color" abstract:"SpecialColors.CursorText"`
	CursorGuide  *Color `plist:"Cursor Guide Color" abstract:"SpecialColors.FindMatch"`
	Link         *Color `plist:"Link Color" abstract:"SpecialColors.Links"`
	SelectedText *Color `plist:"Selected Text Color" abstract:"SpecialColors.SelectedText"`
	Selection    *Color `plist:"Selection Color" abstract:"SpecialColors.Selection"`
}

func (rw *ItermScheme) Name() string {
	return "iterm"
}

func (rw *ItermScheme) TemplateName() string {
	return "iterm.itermcolors.tmpl"
}

// FromString deserializes the plist data into the ItermScheme
func (rw *ItermScheme) FromString(s string) error {
	decoder := plist.NewDecoder(bytes.NewReader([]byte(s)))
	return decoder.Decode(rw)
}
