package goflare

import (
	"fmt"

	"github.com/cdvelop/tinywasm"
)

type Goflare struct {
	tw               *tinywasm.TinyWasm
	config           *Config
	outputJsFileName string // e.g., "_worker.js"
}

type Config struct {
	AppRootDir                string // default: "."
	RelativeInputDirectory    string // input relative directory for source code server app.go to deploy app.wasm (relative) default: "web"
	RelativeOutputDirectory   string // output relative directory for worker.js and app.wasm file (relative) default: "deploy/cloudflare"
	MainInputFile             string // eg: "main.go"
	Logger                    func(message ...any)
	CompilingArguments        func() []string
	OutputWasmFileName        string // WASM file name (default: "worker.wasm")
	BuildPageFunctionShortcut string // build assets wasm,js, json files to pages functions (default: "f")
	BuildWorkerShortcut       string // build assets wasm,js, json files to workers (default: "w")
}

// DefaultConfig returns a Config with all default values set
// AppRootDir=".", RelativeInputDirectory="web", RelativeOutputDirectory="deploy/cloudflare", MainInputFile="main.worker.go", OutputWasmFileName="worker.wasm"
func DefaultConfig() *Config {
	return &Config{
		AppRootDir:              ".",
		RelativeInputDirectory:  "web",
		RelativeOutputDirectory: "deploy/cloudflare",
		MainInputFile:           "main.go",
		Logger:                  func(message ...any) { fmt.Println(message...) },
		CompilingArguments:      nil,
		OutputWasmFileName:      "worker.wasm",

		BuildPageFunctionShortcut: "f",
		BuildWorkerShortcut:       "w",
	}
}

// New creates a new Goflare instance with the provided configuration
// Timeout is set to 40 seconds maximum as TinyGo compilation can be slow
// Default values: AppRootDir=".", RelativeOutputDirectory="deploy/cloudflare", MainInputFile="main.worker.go", OutputWasmFileName="app.wasm"
func New(c *Config) *Goflare {

	dc := DefaultConfig()

	if c == nil {
		c = dc
	} else {
		// Set defaults for empty fields
		if c.AppRootDir == "" {
			c.AppRootDir = dc.AppRootDir
		}

		if c.RelativeInputDirectory == "" {
			c.RelativeInputDirectory = dc.RelativeInputDirectory
		}
		if c.RelativeOutputDirectory == "" {
			c.RelativeOutputDirectory = dc.RelativeOutputDirectory
		}
		if c.MainInputFile == "" {
			c.MainInputFile = dc.MainInputFile
		}
		if c.OutputWasmFileName == "" {
			c.OutputWasmFileName = dc.OutputWasmFileName
		}

		if c.Logger == nil {
			c.Logger = dc.Logger
		}

		if c.BuildPageFunctionShortcut == "" {
			c.BuildPageFunctionShortcut = dc.BuildPageFunctionShortcut
		}
		if c.BuildWorkerShortcut == "" {
			c.BuildWorkerShortcut = dc.BuildWorkerShortcut
		}
	}

	// Extract output name from OutputWasmFileName (remove .wasm extension)
	outputName := c.OutputWasmFileName
	if len(outputName) > 5 && outputName[len(outputName)-5:] == ".wasm" {
		outputName = outputName[:len(outputName)-5]
	}

	tw := tinywasm.New(&tinywasm.Config{
		AppRootDir:              c.AppRootDir,
		SourceDir:               c.RelativeInputDirectory,
		OutputDir:               c.RelativeOutputDirectory,
		WasmExecJsOutputDir:     c.RelativeOutputDirectory,
		MainInputFile:           c.MainInputFile,
		OutputName:              outputName,
		Logger:                  c.Logger,
		CompilingArguments:      c.CompilingArguments,
		DisableWasmExecJsOutput: true, // Pages Advanced Mode embeds wasm_exec.js inline
	})

	g := &Goflare{
		tw:               tw,
		config:           c,
		outputJsFileName: "_worker.js",
	}

	return g
}

// SetCompilerMode changes the compiler mode
// mode: "L" (Large fast/Go), "M" (Medium TinyGo debug), "S" (Small TinyGo production)
func (g *Goflare) SetCompilerMode(newValue string, progress chan<- string) {
	// Execute mode change
	g.tw.Change(newValue, progress)

}
