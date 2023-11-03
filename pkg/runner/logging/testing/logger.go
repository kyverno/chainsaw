package testing

import (
	"github.com/fatih/color"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Logger interface {
	Log(string, *color.Color, ...interface{})
	WithResource(ctrlclient.Object) Logger
}
