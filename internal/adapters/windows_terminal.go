package adapters

import (
	"encoding/json"
)

type WindowsTerminalScheme struct {
	Name                OptString `json:"name" abstract:"Metadata.Name"`
	Black               OptColor  `json:"black" abstract:"AnsiColors.Black"`
	Red                 OptColor  `json:"red" abstract:"AnsiColors.Red"`
	Green               OptColor  `json:"green" abstract:"AnsiColors.Green"`
	Yellow              OptColor  `json:"yellow" abstract:"AnsiColors.Yellow"`
	Blue                OptColor  `json:"blue" abstract:"AnsiColors.Blue"`
	Purple              OptColor  `json:"purple" abstract:"AnsiColors.Magenta"`
	Cyan                OptColor  `json:"cyan" abstract:"AnsiColors.Cyan"`
	White               OptColor  `json:"white" abstract:"AnsiColors.White"`
	BrightBlack         OptColor  `json:"brightBlack" abstract:"AnsiColors.BrightBlack"`
	BrightRed           OptColor  `json:"brightRed" abstract:"AnsiColors.BrightRed"`
	BrightGreen         OptColor  `json:"brightGreen" abstract:"AnsiColors.BrightGreen"`
	BrightYellow        OptColor  `json:"brightYellow" abstract:"AnsiColors.BrightYellow"`
	BrightBlue          OptColor  `json:"brightBlue" abstract:"AnsiColors.BrightBlue"`
	BrightPurple        OptColor  `json:"brightPurple" abstract:"AnsiColors.BrightMagenta"`
	BrightCyan          OptColor  `json:"brightCyan" abstract:"AnsiColors.BrightCyan"`
	BrightWhite         OptColor  `json:"brightWhite" abstract:"AnsiColors.BrightWhite"`
	Background          OptColor  `json:"background" abstract:"SpecialColors.Background"`
	Foreground          OptColor  `json:"foreground" abstract:"SpecialColors.Foreground"`
	SelectionBackground OptColor  `json:"selectionBackground" abstract:"SpecialColors.Selection"`
	CursorColor         OptColor  `json:"cursorColor" abstract:"SpecialColors.Cursor"`
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
