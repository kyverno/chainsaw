package runner

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"k8s.io/utils/clock"
)

type Context struct {
	clock         clock.PassiveClock
	config        v1alpha1.ConfigurationSpec
	ctx           context.Context //nolint:containedctx
	clientFactory func(*testing.T, logging.Logger) client.Client
	namespacer    namespacer.Namespacer
}
