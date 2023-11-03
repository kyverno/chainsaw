package logging

type TLogger interface {
	Log(args ...interface{})
	Helper()
}
