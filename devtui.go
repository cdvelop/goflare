package goflare

func (h *Goflare) Name() string { return "goflare" }
func (h *Goflare) Label() string {
	return "Build Workers"

}
func (h *Goflare) Value() string {
	return ""
}
func (h *Goflare) Change(newValue string, progress func(msgs ...any)) {
	var err error

	switch newValue {
	case h.config.BuildPageFunctionShortcut:
		if progress != nil {
			progress("Starting Pages build...")
		}
		err = h.GeneratePagesFiles()
		if err != nil {
			if progress != nil {
				progress("Pages build failed:", err.Error())
			}
			return
		}
		if progress != nil {
			progress("Pages build completed successfully")
		}

	case h.config.BuildWorkerShortcut:
		if progress != nil {
			progress("Starting Workers build...")
		}
		err = h.GenerateWorkerFiles()
		if err != nil {
			if progress != nil {
				progress("Workers build failed:", err.Error())
			}
			return
		}
		if progress != nil {
			progress("Workers build completed successfully")
		}

	default:
		if progress != nil {
			progress("Unknown shortcut:", newValue)
		}
	}
}

func (h *Goflare) Shortcuts() map[string]string {
	return map[string]string{
		h.config.BuildPageFunctionShortcut: "Build Cloudflare Pages Files",
		h.config.BuildWorkerShortcut:       "Build Cloudflare Workers Files",
	}
}
