package runner

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
)

type Context struct {
	config        v1alpha1.ConfigurationSpec
	ctx           context.Context //nolint:containedctx
	clientFactory func(*testing.T, logging.Logger) client.Client
	namespacer    namespacer.Namespacer
}
