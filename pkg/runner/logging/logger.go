package logging

import (
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Logger interface {
	Log(...interface{})
	Logf(string, ...interface{})
	WithName(string) Logger
	WithResource(ctrlclient.ObjectKey, ctrlclient.Object) Logger
}
