package testing

import (
	"github.com/fatih/color"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Operation string

type Logger interface {
	Log(Operation, *color.Color, ...interface{})
	WithResource(ctrlclient.Object) Logger
}
