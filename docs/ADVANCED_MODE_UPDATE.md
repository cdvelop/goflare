# Pages Advanced Mode Update Summary

## What Changed?

The documentation has been updated to use **Cloudflare Pages Advanced Mode** instead of the standard file-based routing approach.

## Key Differences

### Before (Standard Mode - File-based Routing)
```
functions/
└── api/
    └── [[routes]].mjs        # Catch-all route handler
build/
├── app.wasm
├── wasm_exec.js
├── runtime.mjs
└── worker.mjs
pages/
└── index.html
```

### After (Advanced Mode - Single Worker)
```
pages/
├── _worker.js                # Single entry point (Advanced Mode)
├── app.wasm
├── wasm_exec.js
├── runtime.mjs
└── index.html
```

## Why Advanced Mode?

### Advantages

1. **Simpler Structure** ✅
   - All files in one directory
   - No `/functions` subdirectories
   - Flat, easy-to-understand layout

2. **Full Control** ✅
   - Complete control over routing
   - Custom logic in one place
   - Manual request handling

3. **Worker Compatibility** ✅
   - Same syntax as Workers
   - Easy migration between Workers/Pages
   - Consistent development experience

4. **Better for Go/WASM** ✅
   - No need to split routing between JS files
   - All routing logic in Go code
   - Simpler deployment

5. **Static Assets via Binding** ✅
   - Use `env.ASSETS.fetch()` to serve files
   - Programmatic control over static files
   - No magic file-based routing

## How It Works

### The `_worker.js` File

```javascript
export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    
    // Route 1: API requests → Go WASM
    if (url.pathname.startsWith("/api/")) {
      // Initialize WASM and handle request
      return binding.handleRequest(request);
    }
    
    // Route 2: Everything else → Static assets
    // CRITICAL: Required for serving HTML/CSS/JS
    return env.ASSETS.fetch(request);
  }
};
```

### The `env.ASSETS` Binding

- **Purpose**: Serves static files from your `pages/` directory
- **Usage**: `env.ASSETS.fetch(request)`
- **Automatic**: Cloudflare provides this binding in Pages
- **Required**: Without it, static files won't be served

## Implementation Impact

### GeneratePagesFiles() Method

**Changes Required:**
1. Generate single `_worker.js` file (not `[[routes]].mjs`)
2. Output everything to `pages/` directory (not `build/` + `functions/`)
3. Include `env.ASSETS.fetch()` fallback in template
4. Use `.js` extension (not `.mjs`) for `_worker.js`

**Simplified Steps:**
```go
func (g *Goflare) GeneratePagesFiles() error {
    // 1. Create pages/ directory
    // 2. Copy wasm_exec.js → pages/
    // 3. Generate runtime.mjs → pages/
    // 4. Generate _worker.js → pages/
    return nil
}
```

### Asset Templates

**Before:**
```
assets/
├── runtime/cloudflare.mjs
├── common/worker.mjs
└── pages/routes.mjs      # Catch-all route template
```

**After:**
```
assets/
├── runtime/cloudflare.mjs
├── common/worker.mjs     # For Workers
└── pages/worker.js       # For Pages (Advanced Mode)
```

## Documentation Updates

### BUILD_PAGES.md
- ✅ Updated to Advanced Mode architecture
- ✅ Single `_worker.js` entry point
- ✅ Removed `/functions` directory references
- ✅ Added `env.ASSETS.fetch()` explanation
- ✅ Simplified file structure
- ✅ Updated examples and templates

### ARCHITECTURE.md
- ✅ Updated comparison table
- ✅ Simplified Pages output structure
- ✅ Updated worker.mjs template for Pages
- ✅ Removed `[[routes]].mjs` references

### REVIEW.md
- ✅ Updated output structure
- ✅ Updated differences table
- ✅ Updated asset requirements
- ✅ Added new questions for review

## Migration Notes

### If You Were Using Standard Mode

**Old Approach:**
- Multiple directories (`functions/`, `build/`, `pages/`)
- File-based routing with `[[routes]].mjs`
- Named export `onRequest`

**New Approach (Advanced Mode):**
- Single directory (`pages/`)
- Single `_worker.js` file
- Default export with `fetch()` handler
- `env.ASSETS.fetch()` for static files

### Configuration Changes

**Before:**
```toml
name = "my-app"
pages_build_output_dir = "pages"
```

**After:**
```toml
name = "my-app"
# No pages_build_output_dir needed
# wrangler auto-detects _worker.js in deployment directory
```

## Examples

### Example 1: API + Static Site

```javascript
// pages/_worker.js
export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    
    if (url.pathname.startsWith("/api/")) {
      // Your Go WASM API handlers
      return handleAPI(request, env, ctx);
    }
    
    // Serve index.html, style.css, etc.
    return env.ASSETS.fetch(request);
  }
};
```

### Example 2: Custom 404 Page

```javascript
export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    
    if (url.pathname.startsWith("/api/")) {
      return handleAPI(request, env, ctx);
    }
    
    // Try to serve static asset
    const response = await env.ASSETS.fetch(request);
    
    // Custom 404 handling
    if (response.status === 404) {
      return env.ASSETS.fetch("/404.html");
    }
    
    return response;
  }
};
```

## References

- [Cloudflare Pages Advanced Mode](https://developers.cloudflare.com/pages/functions/advanced-mode/)
- [Workers Static Assets Binding](https://developers.cloudflare.com/workers/static-assets/binding/)
- [env.ASSETS.fetch() API](https://developers.cloudflare.com/pages/functions/api-reference/)

## Next Steps

1. ✅ Review updated documentation
2. ⏳ Approve Advanced Mode approach
3. ⏳ Implement `GeneratePagesFiles()` with new structure
4. ⏳ Create example project using Advanced Mode
5. ⏳ Test deployment to Cloudflare Pages

---

**Status**: Documentation updated, awaiting review ⏳
