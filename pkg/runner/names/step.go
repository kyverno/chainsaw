package names

import (
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func Step(step v1alpha1.TestStep, i int) string {
	if step.Name != "" {
		return step.Name
	}
	return fmt.Sprintf("step-%d", i+1)
}
