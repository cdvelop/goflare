package goflare

import (
	"path"

	"github.com/cdvelop/tinywasm"
)

type Goflare struct {
	tw *tinywasm.TinyWasm
}

type Config struct {
	AppRootDir                 string // default: "."
	WorkerDirSubRelativeOutput string // output path for worker.js and app.wasm file (relative) eg: "deploy"
	MainInputFile              string // eg: "main.worker.go"
	Logger                     func(message ...any)
	CompilingArguments         func() []string
}

// New creates a new Goflare instance with the provided configuration
// Timeout is set to 40 seconds maximum as TinyGo compilation can be slow
// Default values: mainInputFile="main.wasm.go"

func New(c *Config) *Goflare {

	outputFilesDir := path.Join(c.AppRootDir, c.WorkerDirSubRelativeOutput)

	tw := tinywasm.New(&tinywasm.Config{
		AppRootDir:                  c.AppRootDir,
		WebFilesRootRelative:        outputFilesDir,
		WebFilesSubRelative:         outputFilesDir,
		WebFilesSubRelativeJsOutput: outputFilesDir,
		MainInputFile:               c.MainInputFile,
		Logger:                      c.Logger,
		CompilingArguments:          c.CompilingArguments,
	})

	g := &Goflare{
		tw: tw,
	}

	return g
}
