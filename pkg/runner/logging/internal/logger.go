package internal

type Logger interface {
	Log(args ...interface{})
	Helper()
}
