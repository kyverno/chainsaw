package logging

type Operation string

const (
	Apply    Operation = "APPLY"
	Assert   Operation = "ASSERT"
	Catch    Operation = "CATCH"
	Cleanup  Operation = "CLEANUP"
	Command  Operation = "CMD"
	Create   Operation = "CREATE"
	Delete   Operation = "DELETE"
	Error    Operation = "ERROR"
	Finally  Operation = "FINALLY"
	Get      Operation = "GET"
	Internal Operation = "INTERNAL"
	Patch    Operation = "PATCH"
	Script   Operation = "SCRIPT"
	Sleep    Operation = "SLEEP"
	Stderr   Operation = "STDERR"
	Stdout   Operation = "STDOUT"
	Try      Operation = "TRY"
	Update   Operation = "UPDATE"
)
