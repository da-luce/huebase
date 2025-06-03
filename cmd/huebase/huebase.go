package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/da-luce/huebase/internal/adapter"
	log "github.com/sirupsen/logrus"
)

var schemes = map[string]adapter.Adapter{}

func init() {
	for _, adapter := range adapter.Adapters {
		schemes[adapter.Name()] = adapter
	}
}

func main() {
	// Define CLI arguments
	inputFile := flag.String("input", "", "Path to the input theme file")
	outputFile := flag.String("output", "", "Path to save the converted theme file (optional)")
	inputFormat := flag.String("from", "", "Input format (e.g., base16)")
	outputFormat := flag.String("to", "", "Output format (e.g., alacritty)")

	flag.Parse()

	log.SetLevel(log.InfoLevel)

	// Log the received CLI arguments
	log.WithFields(log.Fields{
		"inputFile":    *inputFile,
		"outputFile":   *outputFile,
		"inputFormat":  *inputFormat,
		"outputFormat": *outputFormat,
	}).Info("CLI arguments parsed")

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

	data, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Failed to read input file: %v", err)
		os.Exit(1)
	}

	// Perform the conversion
	outputData, err := adapter.ConvertTheme(string(data), reader, writer)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Write to output file if specified, else print to stdout
	if *outputFile != "" {
		err = os.WriteFile(*outputFile, []byte(outputData), 0644)
		if err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Converted theme saved to %s\n", *outputFile)
	} else {
		fmt.Println(outputData)
	}
}
