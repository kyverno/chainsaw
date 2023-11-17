package logging

import (
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
)

type (
	Logger    = tlogging.Logger
	Operation = tlogging.Operation
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
	Std___   Operation = "STD___"
	Stderr   Operation = "STDERR"
	Stdout   Operation = "STDOUT"
	Try      Operation = "TRY"
)
