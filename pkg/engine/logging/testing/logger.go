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
	WithCluster(*string) Logger
	WithResource(client.Object) Logger
}
