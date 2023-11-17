package processors

import (
	"time"

	"github.com/kyverno/chainsaw/pkg/runner/operations"
)

type operation struct {
	continueOnError bool
	timeout         time.Duration
	operation       operations.Operation
}
