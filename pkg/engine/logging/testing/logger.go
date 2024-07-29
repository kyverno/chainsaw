package testing

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
)

type (
	Operation string
	Status    string
)

type Logger interface {
	Log(Operation, Status, *color.Color, ...fmt.Stringer)
	WithResource(client.Object) Logger
}
