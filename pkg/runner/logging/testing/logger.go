package testing

import (
	"fmt"

	"github.com/fatih/color"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	Operation string
	Status    string
)

type Logger interface {
	Log(Operation, Status, *color.Color, ...fmt.Stringer)
	WithResource(ctrlclient.Object) Logger
}
