package logging

import (
	"fmt"

	"github.com/fatih/color"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type MockLogger struct {
	Logs []string
}

func (m *MockLogger) WithResource(resource ctrlclient.Object) Logger {
	return m
}

func (m *MockLogger) Log(operation string, color *color.Color, args ...interface{}) {
	message := fmt.Sprintf("%s: %v", operation, args)
	m.Logs = append(m.Logs, message)
}
