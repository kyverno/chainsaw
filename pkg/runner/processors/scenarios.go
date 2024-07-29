package processors

import (
	"github.com/kyverno/chainsaw/pkg/discovery"
)

func applyScenarios(test discovery.Test) []discovery.Test {
	var scenarios []discovery.Test
	if test.Test != nil {
		if len(test.Test.Spec.Scenarios) == 0 {
			scenarios = append(scenarios, test)
		} else {
			for s := range test.Test.Spec.Scenarios {
				scenario := test.Test.Spec.Scenarios[s]
				test := test
				test.Test = test.Test.DeepCopy()
				test.Test.Spec.Scenarios = nil
				bindings := scenario.Bindings
				bindings = append(bindings, test.Test.Spec.Bindings...)
				test.Test.Spec.Bindings = bindings
				scenarios = append(scenarios, test)
			}
		}
	}
	return scenarios
}
