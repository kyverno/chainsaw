package runner

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/testing"
)

func runTest(goctx context.Context, ctx Context, config v1alpha1.ConfigurationSpec, test discovery.Test, size int) {
	t := testing.FromContext(goctx)
	t.Helper()
	for i, step := range test.Spec.Steps {
		name := step.Name
		if name == "" {
			name = fmt.Sprintf("step-%d", i+1)
		}
		goctx := logging.IntoContext(goctx, logging.NewLogger(t, ctx.clock, test.Name, fmt.Sprintf("%-*s", size, name)))
		executeStep(goctx, ctx, test.BasePath, config, test.Spec, step)
	}
}
