package logging

type Status string

const (
	BeginStatus Status = "BEGIN"
	DoneStatus  Status = "DONE"
	EndStatus   Status = "END"
	ErrorStatus Status = "ERROR"
	LogStatus   Status = "LOG"
	OkStatus    Status = "OK"
	RunStatus   Status = "RUN"
	WarnStatus  Status = "WARN"
	// DoneStatus  Status = "✅"
	// ErrorStatus Status = "❌"
	// LogStatus   Status = "📄"
	// OkStatus    Status = "🟢"
	// RunStatus   Status = "🚧"
	// WarnStatus  Status = "🟡"
)
