package check

import (
	"context"
	"errors"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Check(ctx context.Context, obj interface{}, bindings binding.Bindings, check *v1alpha1.Check) (field.ErrorList, error) {
	if check == nil {
		return nil, errors.New("check is null")
	}
	if check.Value == nil {
		return nil, errors.New("check value is null")
	}
	return assert.Assert(ctx, assert.Parse(ctx, check.Value), obj, bindings)
}
