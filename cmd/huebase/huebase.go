package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/da-luce/huebase/internal/adapters"
)

var schemes = map[string]adapters.Adapter{
	"base16":           &adapters.Base16Scheme{},
	"16":               &adapters.Base16Scheme{},
	"alacritty":        &adapters.AlacrittyScheme{},
	"windows_terminal": &adapters.WindowsTerminalScheme{},
	"wt":               &adapters.WindowsTerminalScheme{},
	"gogh":             &adapters.GoghScheme{},
}

func convertTheme(inputFile string, inputFormat string, outputFormat string) (string, error) {
	// Ensure the formats are supported
	inputReader, ok := schemes[inputFormat]
	if !ok {
		return "", errors.New("unsupported input format: " + inputFormat)
	}
	outputWriter, ok := schemes[outputFormat]
	if !ok {
		return "", errors.New("unsupported output format: " + outputFormat)
	}

	// Read the input file
	inputData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read input file: %v", err)
	}

	// Convert to struct (populate inputReader directly)
	if err := inputReader.FromString(string(inputData)); err != nil {
		return "", fmt.Errorf("failed to parse input file: %v", err)
	}

	// Convert to AbstractTheme
	abstractTheme, err := adapters.ToAbstract(inputReader)
	if err != nil {
		return "", fmt.Errorf("failed to convert to AbstractTheme: %v", err)
	}

	// Convert AbstractTheme to output format
	adapters.FromAbstract(&abstractTheme, outputWriter)

	// Convert output to string
	outputData, err := outputWriter.ToString()
	if err != nil {
		return "", fmt.Errorf("failed to convert theme to output format: %v", err)
	}

	return outputData, nil
}

func main() {
	// Define CLI arguments
	inputFile := flag.String("input", "", "Path to the input theme file")
	outputFile := flag.String("output", "", "Path to save the converted theme file (optional)")
	inputFormat := flag.String("from", "", "Input format (e.g., base16)")
	outputFormat := flag.String("to", "", "Output format (e.g., alacritty)")

	flag.Parse()

	// Validate required arguments
	if *inputFile == "" || *inputFormat == "" || *outputFormat == "" {
		fmt.Println("Error: Missing required arguments")
		flag.Usage()
		os.Exit(1)
	}

	// Perform the conversion
	outputData, err := convertTheme(*inputFile, *inputFormat, *outputFormat)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Write to output file if specified, else print to stdout
	if *outputFile != "" {
		err = ioutil.WriteFile(*outputFile, []byte(outputData), 0644)
		if err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Converted theme saved to %s\n", *outputFile)
	} else {
		fmt.Println(outputData)
	}
}
