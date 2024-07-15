package discovery

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

type Test struct {
	Test     *v1alpha1.Test
	BasePath string
	Err      error
}
