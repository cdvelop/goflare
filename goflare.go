package goflare

import "github.com/cdvelop/tinywasm"

type Goflare struct {
	tw *tinywasm.TinyWasm
}

type Config struct {
	AppRootDir                  string
	WebFilesRootRelative        string
	WebFilesSubRelative         string
	WebFilesSubRelativeJsOutput string
	Logger                      func(message ...any)
	CodingShortcut              string
	DebuggingShortcut           string
	ProductionShortcut          string
	Callback                    func(err error)
	CompilingArguments          func() []string
}

// New creates a new Goflare instance with the provided configuration
// Timeout is set to 40 seconds maximum as TinyGo compilation can be slow
// Default values: mainInputFile="main.wasm.go"

func New() *Goflare {

	tw := tinywasm.New(&tinywasm.Config{
		AppRootDir:                  ".",
		WebFilesRootRelative:        "web",
		WebFilesSubRelative:         "public",
		WebFilesSubRelativeJsOutput: "theme/js",
		Logger:                      func(message ...any) {},
		CodingShortcut:              "c",
		DebuggingShortcut:           "d",
		ProductionShortcut:          "p",
		Callback:                    func(err error) {},
		CompilingArguments:          func() []string { return []string{} },
	})

	g := &Goflare{
		tw: tw,
	}

	return g
}
