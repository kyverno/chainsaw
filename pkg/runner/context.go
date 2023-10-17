package runner

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
)

type Context struct {
	clientFactory func(*testing.T, logging.Logger) client.Client
	namespacer    namespacer.Namespacer

	passedTests *int
	failedTests *int
}
