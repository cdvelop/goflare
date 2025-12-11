# Goflare CLI

Command-line tool to generate Cloudflare Pages files from Go WASM projects.

## Usage

Run this command from your project root directory:

```bash
goflare
```

## Requirements

Your project must have:
- A Go WASM source file at: `web/main.worker.go`
- TinyGo installed (for optimal WASM compilation)

## Project Structure

```
your-project/
├── go.mod              # Your project's go.mod
├── web/
│   └── main.worker.go  # Your WASM entry point
└── deploy/
    └── cloudflare/     # Generated files will be here
        ├── _worker.js
        └── worker.wasm
```

## Example main.worker.go

```go
package main

import "syscall/js"

func main() {
    js.Global().Set("myFunction", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        return "Hello from WASM!"
    }))
    select {}
}
```

## Installation

```bash
go install github.com/tinywasm/goflare/cmd/goflare@latest
```

## Output

The command generates two files:
- `deploy/cloudflare/_worker.js` - Cloudflare Pages Advanced Mode worker
- `deploy/cloudflare/worker.wasm` - Compiled WASM binary
