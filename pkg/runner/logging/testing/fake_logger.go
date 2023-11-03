package testing

import (
	"fmt"

	"github.com/fatih/color"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type FakeLogger struct {
	Logs []string
}

func (m *FakeLogger) WithResource(resource ctrlclient.Object) Logger {
	return m
}

func (m *FakeLogger) Log(operation string, color *color.Color, args ...interface{}) {
	message := fmt.Sprintf("%s: %v", operation, args)
	m.Logs = append(m.Logs, message)
}
