package operations

// import (
// 	"fmt"

// 	"github.com/fatih/color"
// 	"github.com/kyverno/chainsaw/pkg/runner/logging"
// 	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
// )

// type MockLogger struct {
// 	Logs []string
// }

// func (m *MockLogger) WithResource(resource ctrlclient.Object) logging.Logger {
// 	return m
// }

// func (m *MockLogger) Log(operation string, color *color.Color, args ...interface{}) {
// 	message := fmt.Sprintf("%s: %v", operation, args)
// 	m.Logs = append(m.Logs, message)
// }
