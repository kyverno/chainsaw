package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/engine/bindings"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
)

func RegisterWarningsInBindings(ctx context.Context, tc apis.Bindings) apis.Bindings {
	ctxTc := enginecontext.TestContextFromCtx(ctx)
	if ctxTc != nil && ctxTc.CurrentCluster() != nil {
		warnings := ctxTc.CurrentCluster().GetWarnings()
		tc = bindings.RegisterBinding(tc, "warnings", warnings)
	}

	return tc
}

func ResetWarnings(ctx context.Context) {
	ctxTc := enginecontext.TestContextFromCtx(ctx)
	if ctxTc != nil && ctxTc.CurrentCluster() != nil {
		ctxTc.CurrentCluster().ResetWarnings()
	}
}
