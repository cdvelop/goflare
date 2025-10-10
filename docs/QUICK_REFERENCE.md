# GoFlare Quick Reference

## Workers vs Pages (Advanced Mode)

### Workers Deployment

**Output Directory**: `deploy/`

**Files Generated**:
- `worker.mjs` - Main entry point
- `app.wasm` - Compiled Go binary
- `wasm_exec.js` - WASM runtime
- `runtime.mjs` - Cloudflare context

**Entry Point**:
```javascript
export default {
  fetch(request, env, ctx) { /* ... */ }
};
```

**Deploy**:
```bash
wrangler deploy
```

**Config** (`wrangler.toml`):
```toml
name = "my-worker"
main = "./deploy/worker.mjs"
compatibility_date = "2023-04-30"
```

---

### Pages Deployment (Advanced Mode)

**Output Directory**: `pages/`

**Files Generated**:
- `_worker.js` - Advanced Mode entry point
- `app.wasm` - Compiled Go binary
- `wasm_exec.js` - WASM runtime
- `runtime.mjs` - Cloudflare context
- Plus your static files (HTML, CSS, JS)

**Entry Point**:
```javascript
export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    
    // API routes → Go WASM
    if (url.pathname.startsWith("/api/")) {
      return handleAPI(request, env, ctx);
    }
    
    // Everything else → Static files
    return env.ASSETS.fetch(request);
  }
};
```

**Deploy**:
```bash
wrangler pages deploy ./pages
```

**Config** (`wrangler.toml`):
```toml
name = "my-pages-app"
compatibility_date = "2023-04-30"
# No other config needed - auto-detects _worker.js
```

---

## Key Differences

| Feature | Workers | Pages (Advanced Mode) |
|---------|---------|---------------------|
| **File name** | `worker.mjs` | `_worker.js` |
| **Extension** | `.mjs` | `.js` |
| **Static files** | ❌ Not supported | ✅ Via `env.ASSETS.fetch()` |
| **Directory** | `deploy/` | `pages/` |
| **Command** | `wrangler deploy` | `wrangler pages deploy ./pages` |
| **Config field** | `main = ...` | Auto-detected |

---

## GoFlare Usage

### Setup

```go
import "github.com/cdvelop/goflare"

g := goflare.New(&goflare.Config{
    AppRootDir:                 ".",
    RelativeOutputDirectory: "deploy", // or "pages"
    MainInputFile:              "main.go",
    Logger: func(message ...any) {
        log.Println(message...)
    },
})
```

### Generate Workers Files

```go
if err := g.GenerateWorkerFiles(); err != nil {
    log.Fatal(err)
}
```

**Outputs**: `deploy/worker.mjs`, `deploy/app.wasm`, etc.

### Generate Pages Files

```go
if err := g.GeneratePagesFiles(); err != nil {
    log.Fatal(err)
}
```

**Outputs**: `pages/_worker.js`, `pages/app.wasm`, etc.

---

## Build Script Example

```go
// main.builder.go
package main

import (
    "log"
    "github.com/cdvelop/goflare"
)

func main() {
    g := goflare.New(&goflare.Config{
        AppRootDir:                 ".",
        RelativeOutputDirectory: "pages",
        MainInputFile:              "main.go",
        Logger:                     log.Println,
    })
    
    // For Pages
    if err := g.GeneratePagesFiles(); err != nil {
        log.Fatal(err)
    }
    
    // OR for Workers
    // if err := g.GenerateWorkerFiles(); err != nil {
    //     log.Fatal(err)
    // }
}
```

**Run**:
```bash
go run main.builder.go
```

---

## Makefile Example

```makefile
.PHONY: build-pages
build-pages:
	go run main.builder.go

.PHONY: build-workers
build-workers:
	go run main.builder.go

.PHONY: dev-pages
dev-pages: build-pages
	wrangler pages dev ./pages

.PHONY: dev-workers
dev-workers: build-workers
	wrangler dev

.PHONY: deploy-pages
deploy-pages: build-pages
	wrangler pages deploy ./pages

.PHONY: deploy-workers
deploy-workers: build-workers
	wrangler deploy
```

---

## Project Structure

### Pages Project

```
my-pages-app/
├── main.go              # Your Go application
├── main.builder.go      # Build script
├── wrangler.toml
├── go.mod
└── pages/               # Generated (deploy this)
    ├── _worker.js
    ├── app.wasm
    ├── wasm_exec.js
    ├── runtime.mjs
    ├── index.html       # Your static files
    └── style.css
```

### Workers Project

```
my-worker/
├── main.go              # Your Go application
├── main.builder.go      # Build script
├── wrangler.toml
├── go.mod
└── deploy/              # Generated (deploy this)
    ├── worker.mjs
    ├── app.wasm
    ├── wasm_exec.js
    └── runtime.mjs
```

---

## Common Patterns

### API Routes

```go
// In your Go code (main.go)
package main

import "github.com/syumai/workers"

func main() {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasPrefix(r.URL.Path, "/api/") {
            // Handle API requests
            json.NewEncoder(w).Encode(map[string]string{
                "message": "Hello from Go!",
            })
            return
        }
        // For Pages, other routes are handled by env.ASSETS.fetch()
    })
    
    workers.Serve(handler)
}
```

### Static + Dynamic

For Pages, your `_worker.js` routes:
- `/api/*` → Go WASM (dynamic)
- `/*` → `env.ASSETS.fetch()` (static)

---

## Troubleshooting

### Static files not loading (Pages)
✅ Make sure `_worker.js` includes `return env.ASSETS.fetch(request)`

### Worker not found
✅ Check file name: `worker.mjs` (Workers) vs `_worker.js` (Pages)

### Module import errors
✅ Verify relative imports: `import "./wasm_exec.js"`

### WASM not initializing
✅ Check that `wasm_exec.js` matches your compilation mode (Go vs TinyGo)

---

## Resources

- [BUILD_WORKERS.md](BUILD_WORKERS.md) - Detailed Workers specification
- [BUILD_PAGES.md](BUILD_PAGES.md) - Detailed Pages specification
- [ARCHITECTURE.md](ARCHITECTURE.md) - Architecture overview
- [Cloudflare Docs](https://developers.cloudflare.com/pages/functions/advanced-mode/)
