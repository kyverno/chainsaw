package discovery

import (
	"github.com/kyverno/chainsaw/pkg/model"
)

type Test struct {
	Test     *model.Test
	BasePath string
	Err      error
}
