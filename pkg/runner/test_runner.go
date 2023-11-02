package runner

import (
	"fmt"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
)

func runTest(t *testing.T, ctx Context, config v1alpha1.ConfigurationSpec, test discovery.Test, size int) {
	t.Helper()
	for i, step := range test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		logger := logging.NewLogger(t, ctx.clock, test.Name, fmt.Sprintf("%-*s", size, name))
		executeStep(t, logger, ctx, test.BasePath, config, test.Spec, step)
	}
}
