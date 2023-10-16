package runner

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
)

type Context struct {
	clientFactory func(*testing.T) client.Client
	namespacer    namespacer.Namespacer
}
