package logging

type TLogger interface {
	Log(args ...any)
	Helper()
}
