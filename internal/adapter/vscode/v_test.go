package vscode_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/da-luce/paletteport/internal/adapter/vscode"
)

func TestLoadThemeFromFile(t *testing.T) {
	filePath := "./tokyo-night-color-theme.json" // replace with your actual file path

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read theme file: %v", err)
	}

	var theme vscode.VSCodeTheme
	err = theme.FromString(string(data))
	if err != nil {
		t.Fatalf("Failed to unmarshal theme JSON: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(theme, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal theme for pretty print: %v", err)
	}
	t.Logf("Parsed theme struct:\n%s", string(prettyJSON))

}
