package logging

import (
	"github.com/kyverno/kyverno/ext/output/color"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Logger interface {
	Log(string, *color.Color, ...interface{})
	WithResource(ctrlclient.Object) Logger
}
