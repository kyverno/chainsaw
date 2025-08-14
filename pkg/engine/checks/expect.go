package checks

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

func Expect(ctx context.Context, compilers compilers.Compilers, obj unstructured.Unstructured, bindings apis.Bindings, expect ...v1alpha1.Expectation) (bool, error) {
	matched := false
	var results field.ErrorList
	for _, expectation := range expect {
		// if a match is specified, skip the check if the resource doesn't match
		if expectation.Match != nil && !expectation.Match.IsNil() {
			if errs, err := Check(ctx, compilers, obj.UnstructuredContent(), nil, expectation.Match); err != nil {
				return true, err
			} else if len(errs) != 0 {
				continue
			}
		}
		matched = true
		if errs, err := Check(ctx, compilers, obj.UnstructuredContent(), bindings, ptr.To(expectation.Check)); err != nil {
			return true, err
		} else {
			results = append(results, errs...)
		}
	}
	return matched, results.ToAggregate()
}
