package model

type WithWarnings struct {
	warnings []map[string]any
}

func NewWithWarnings() *WithWarnings {
	return &WithWarnings{warnings: make([]map[string]any, 0)}
}

type WarningsHolder interface {
	GetWarnings() []map[string]any
	ResetWarnings()
	HandleWarningHeader(code int, agent string, text string)
}

func (w *WithWarnings) GetWarnings() []map[string]any {
	return w.warnings
}

func (w *WithWarnings) ResetWarnings() {
	w.warnings = make([]map[string]any, 0)
}

func (w *WithWarnings) HandleWarningHeader(code int, agent string, text string) {
	w.warnings = append(w.warnings, map[string]any{"code": code, "agent": agent, "text": text})
}
