package runner

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/logging"
	"k8s.io/utils/clock"
)

const eraser = "\b\b\b\b\b\b\b\b\b\b\b\b"

// newSink creates a new sink for logging.
// If quietMode is true, it will buffer non-error messages and flush them when an error occurs.
// This ensures diagnostic details (catch output, failed command output) are shown on failures.
// Otherwise, it will log all messages immediately.
func newSink(clock clock.PassiveClock, quiet bool, log func(args ...any)) logging.SinkFunc {
	// Use a map to maintain separate buffers per test+step combination
	// This ensures parallel tests don't mix their output
	buffers := make(map[string][]string)
	// Track which test+step combinations have encountered errors
	// Once an error occurs, all subsequent logs for that test+step are shown immediately
	errorOccurred := make(map[string]bool)
	
	formatLog := func(test string, step string, operation logging.Operation, status logging.Status, obj client.Object, color *color.Color, args ...fmt.Stringer) string {
		sprint := fmt.Sprint
		opLen := 9
		stLen := 5
		if color != nil {
			sprint = color.Sprint
			opLen += 14
			stLen += 14
		}
		a := make([]any, 0, len(args)+2)
		prefix := fmt.Sprintf("%s| %s | %s | %s | %-*s | %-*s |", eraser, clock.Now().Format("15:04:05"), sprint(test), sprint(step), opLen, sprint(operation), stLen, sprint(status))
		if obj != nil {
			gvk := obj.GetObjectKind().GroupVersionKind()
			key := client.Key(obj)
			prefix = fmt.Sprintf("%s %s/%s @ %s", prefix, gvk.GroupVersion(), gvk.Kind, client.Name(key))
		}
		a = append(a, prefix)
		for _, arg := range args {
			a = append(a, "\n")
			a = append(a, arg)
		}
		return fmt.Sprint(a...)
	}
	
	return func(test string, step string, operation logging.Operation, status logging.Status, obj client.Object, color *color.Color, args ...fmt.Stringer) {
		formatted := formatLog(test, step, operation, status, obj, color, args...)
		bufferKey := test + "|" + step
		
		if !quiet {
			// Not in quiet mode - log everything immediately
			log(formatted)
		} else if status == logging.ErrorStatus || operation == logging.Internal {
			// In quiet mode with error/internal: flush buffer and mark error occurred
			if buffer, exists := buffers[bufferKey]; exists {
				for _, buffered := range buffer {
					log(buffered)
				}
				delete(buffers, bufferKey) // Clear buffer after flushing
			}
			errorOccurred[bufferKey] = true
			log(formatted)
		} else if errorOccurred[bufferKey] {
			// After an error, show all subsequent logs immediately (including catch blocks)
			log(formatted)
		} else {
			// In quiet mode before any error: buffer non-error messages per test+step
			buffers[bufferKey] = append(buffers[bufferKey], formatted)
		}
	}
}
