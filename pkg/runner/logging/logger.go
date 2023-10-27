package logging

import (
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Logger interface {
	Log(string, ...interface{})
	WithResource(ctrlclient.Object) Logger
}
