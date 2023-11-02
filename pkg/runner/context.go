package runner

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"k8s.io/utils/clock"
)

type Context struct {
	clock      clock.PassiveClock
	client     client.Client
	namespacer namespacer.Namespacer
}
