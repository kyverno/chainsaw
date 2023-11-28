package check

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Expectations(ctx context.Context, obj unstructured.Unstructured, bindings binding.Bindings, expect ...v1alpha1.Expectation) (bool, error) {
	matched := false
	var results field.ErrorList
	for _, expectation := range expect {
		// if a match is specified, skip the check if the resource doesn't match
		if expectation.Match != nil && expectation.Match.Value != nil {
			errs, validationErr := assert.Validate(ctx, expectation.Match.Value, obj.UnstructuredContent(), nil)
			if validationErr != nil {
				return true, validationErr
			}
			if len(errs) != 0 {
				continue
			}
		}
		matched = true
		errs, validationErr := assert.Validate(ctx, expectation.Check.Value, obj.UnstructuredContent(), bindings)
		if validationErr != nil {
			return true, validationErr
		}
		results = append(results, errs...)
	}
	return matched, results.ToAggregate()
}
