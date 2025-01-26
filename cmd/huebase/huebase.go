package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/da-luce/huebase/internal/adapters"
)

var schemes = map[string]adapters.Adapter{
	"base16":           &adapters.Base16Scheme{},
	"16":               &adapters.Base16Scheme{},
	"alacritty":        &adapters.AlacrittyScheme{},
	"windows_terminal": &adapters.WindowsTerminalScheme{},
	"wt":               &adapters.WindowsTerminalScheme{},
	"gogh":             &adapters.GoghScheme{},
	"iterm":            &adapters.ItermScheme{},
}

func main() {
	// Define CLI arguments
	inputFile := flag.String("input", "", "Path to the input theme file")
	outputFile := flag.String("output", "", "Path to save the converted theme file (optional)")
	inputFormat := flag.String("from", "", "Input format (e.g., base16)")
	outputFormat := flag.String("to", "", "Output format (e.g., alacritty)")

	flag.Parse()

	log.SetLevel(log.InfoLevel)

	// Validate required arguments
	if *inputFile == "" || *inputFormat == "" || *outputFormat == "" {
		fmt.Println("Error: Missing required arguments")
		flag.Usage()
		os.Exit(1)
	}

	reader, ok := schemes[*inputFormat]
	if !ok {
		fmt.Printf("Unsupported input format: %s\n", *inputFormat)
		os.Exit(1)
	}

	writer, ok := schemes[*outputFormat]
	if !ok {
		fmt.Printf("Unsupported output format: %s,\n", *outputFormat)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Failed to read input file: %v", err)
		os.Exit(1)
	}

	// Perform the conversion
	outputData, err := adapters.ConvertTheme(string(data), reader, writer)
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
