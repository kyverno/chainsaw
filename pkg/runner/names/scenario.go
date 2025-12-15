package names

import (
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func Scenario(scenario v1alpha1.Scenario, i int) string {
	if scenario.Name != "" {
		return scenario.Name
	}
	return fmt.Sprintf("scenario #%d", i+1)
}
