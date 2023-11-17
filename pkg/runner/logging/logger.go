package logging

import (
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
)

type (
	Logger    = tlogging.Logger
	Operation = tlogging.Operation
	Status    = tlogging.Status
)

const (
	Apply    Operation = "APPLY"
	Assert   Operation = "ASSERT"
	Catch    Operation = "CATCH"
	Command  Operation = "CMD"
	Create   Operation = "CREATE"
	Delete   Operation = "DELETE"
	Error    Operation = "ERROR"
	Finally  Operation = "FINALLY"
	Get      Operation = "GET"
	Internal Operation = "INTERNAL"
	Patch    Operation = "PATCH"
	Script   Operation = "SCRIPT"
	Stderr   Operation = "STDERR"
	Stdout   Operation = "STDOUT"
	Try      Operation = "TRY"
)

const (
	DoneStatus  Status = "DONE"
	ErrorStatus Status = "ERROR"
	OkStatus    Status = "OK"
	RunStatus   Status = "RUN"
	LogStatus   Status = "LOG"
)
