package runner

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
)

type Context struct {
	client     client.Client
	namespacer namespacer.Namespacer
}
