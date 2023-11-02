package runner

import (
	"fmt"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/runner/internal"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
)

func Run(cfg *rest.Config, clock clock.PassiveClock, config v1alpha1.ConfigurationSpec, tests ...discovery.Test) (*summary.Summary, error) {
	var summary summary.Summary
	if len(tests) == 0 {
		return &summary, nil
	}
	if err := internal.SetupFlags(config); err != nil {
		return nil, err
	}
	client, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	internalTests := []testing.InternalTest{{
		Name: "chainsaw",
		F: func(t *testing.T) {
			t.Helper()
			runner := &testsRunner{
				config:  config,
				client:  client,
				clock:   clock,
				summary: &summary,
			}
			runner.runTests(t, tests...)
		},
	}}
	deps := &internal.TestDeps{}
	m := testing.MainStart(deps, internalTests, nil, nil, nil)
	if code := m.Run(); code > 1 {
		return &summary, fmt.Errorf("testing framework exited with non zero code %d", code)
	}
	return &summary, nil
}
