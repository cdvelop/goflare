package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cdvelop/goflare"
)

func main() {
	config := goflare.DefaultConfig()

	// Verify input file exists
	inputFile := filepath.Join(config.RelativeInputDirectory, config.MainInputFile)
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Input file not found: %s\n", inputFile)
		fmt.Fprintf(os.Stderr, "Please create your WASM entry point at this location.\n")
		os.Exit(1)
	}

	g := goflare.New(config)
	g.SetCompilerMode("m") // Use TinyGo minimal mode

	fmt.Println("Generating Cloudflare Pages files...")

	if err := g.GeneratePagesFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Files generated successfully!")
	fmt.Printf("  - %s/_worker.js\n", config.RelativeOutputDirectory)
	fmt.Printf("  - %s/%s\n", config.RelativeOutputDirectory, config.OutputWasmFileName)
}
