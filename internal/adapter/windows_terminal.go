package adapter

import (
	"encoding/json"
)

type WindowsTerminalScheme struct {
	Name                *string `json:"name" abstract:"Metadata.Name"`
	Black               *Color  `json:"black" abstract:"AnsiColors.Black"`
	Red                 *Color  `json:"red" abstract:"AnsiColors.Red"`
	Green               *Color  `json:"green" abstract:"AnsiColors.Green"`
	Yellow              *Color  `json:"yellow" abstract:"AnsiColors.Yellow"`
	Blue                *Color  `json:"blue" abstract:"AnsiColors.Blue"`
	Purple              *Color  `json:"purple" abstract:"AnsiColors.Magenta"`
	Cyan                *Color  `json:"cyan" abstract:"AnsiColors.Cyan"`
	White               *Color  `json:"white" abstract:"AnsiColors.White"`
	BrightBlack         *Color  `json:"brightBlack" abstract:"AnsiColors.BrightBlack"`
	BrightRed           *Color  `json:"brightRed" abstract:"AnsiColors.BrightRed"`
	BrightGreen         *Color  `json:"brightGreen" abstract:"AnsiColors.BrightGreen"`
	BrightYellow        *Color  `json:"brightYellow" abstract:"AnsiColors.BrightYellow"`
	BrightBlue          *Color  `json:"brightBlue" abstract:"AnsiColors.BrightBlue"`
	BrightPurple        *Color  `json:"brightPurple" abstract:"AnsiColors.BrightMagenta"`
	BrightCyan          *Color  `json:"brightCyan" abstract:"AnsiColors.BrightCyan"`
	BrightWhite         *Color  `json:"brightWhite" abstract:"AnsiColors.BrightWhite"`
	Background          *Color  `json:"background" abstract:"SpecialColors.Background"`
	Foreground          *Color  `json:"foreground" abstract:"SpecialColors.Foreground"`
	SelectionBackground *Color  `json:"selectionBackground" abstract:"SpecialColors.Selection"`
	CursorColor         *Color  `json:"cursorColor" abstract:"SpecialColors.Cursor"`
}

func (rw *WindowsTerminalScheme) FromString(input string) error {
	err := json.Unmarshal([]byte(input), rw)
	if err != nil {
		return err
	}
	return nil
}

func (rw *WindowsTerminalScheme) ToString() (string, error) {
	data, err := json.MarshalIndent(rw, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
