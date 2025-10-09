package goflare

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GeneratePagesFiles generates all necessary files for Cloudflare Pages Functions using Advanced Mode
func (g *Goflare) GeneratePagesFiles() error {
	// 1. Ensure output directory exists
	pagesDir := filepath.Join(g.tw.Config.AppRootDir, "pages")
	if err := os.MkdirAll(pagesDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create pages directory: %w", err)
	}

	// 2. Update tinywasm configuration to output to pages/ directory with configurable WASM filename
	originalConfig := *g.tw.Config
	g.tw.Config.WebFilesRootRelative = pagesDir
	g.tw.Config.WebFilesSubRelative = pagesDir
	g.tw.Config.WebFilesSubRelativeJsOutput = pagesDir
	// Note: WASM filename configuration would need to be handled by tinywasm

	defer func() {
		// Restore original config
		g.tw.Config.WebFilesRootRelative = originalConfig.WebFilesRootRelative
		g.tw.Config.WebFilesSubRelative = originalConfig.WebFilesSubRelative
		g.tw.Config.WebFilesSubRelativeJsOutput = originalConfig.WebFilesSubRelativeJsOutput
	}()

	// 3. Generate _worker.js (Advanced Mode entry point with everything inline)
	if err := g.generatePagesWorker(pagesDir); err != nil {
		return fmt.Errorf("failed to generate _worker.js: %w", err)
	}

	return nil
}

// generatePagesWorker generates the _worker.js file for Pages (Advanced Mode)
// This creates a single combined file with wasm_exec.js and runtime.mjs inline
func (g *Goflare) generatePagesWorker(pagesDir string) error {
	destPath := filepath.Join(pagesDir, "_worker.js")

	// Read wasm_exec.js content
	wasmExecContent, err := g.getWasmExecContent()
	if err != nil {
		return fmt.Errorf("failed to get wasm_exec content: %w", err)
	}

	// Read and modify runtime.mjs content
	runtimeContent, err := g.getModifiedRuntimeContent()
	if err != nil {
		return fmt.Errorf("failed to get runtime content: %w", err)
	}

	// Read worker template
	workerTemplate, err := assets.ReadFile("assets/pages/worker_combined.js")
	if err != nil {
		return fmt.Errorf("failed to read worker template: %w", err)
	}

	// Combine all content: wasm_exec.js + runtime.mjs + worker logic
	combinedContent := string(wasmExecContent) + "\n\n" + string(runtimeContent) + "\n\n" + string(workerTemplate)

	// Replace placeholders with actual values
	combinedContent = strings.ReplaceAll(combinedContent, "API_ROUTE_PREFIX", g.config.PagesApiRoutePrefix)

	// Write the combined file
	err = os.WriteFile(destPath, []byte(combinedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write combined worker file: %w", err)
	}

	return nil
}

// getModifiedRuntimeContent reads and modifies the runtime.mjs content
func (g *Goflare) getModifiedRuntimeContent() string {

	return fmt.Sprintf(`import { connect } from "cloudflare:sockets";
import mod from "%v";

export async function loadModule() {
  return mod;
}

export function createRuntimeContext({ env, ctx, binding }) {
  return {
    env,
    ctx,
    connect,
    binding,
  };
}`, g.OutFile)
}
