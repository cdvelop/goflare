package goflare

func (h *Goflare) UnobservedFiles() []string {
	return []string{
		h.tw.OutputRelativePath(),
		h.tw.WasmExecJsOutputPath(),
	}
}
