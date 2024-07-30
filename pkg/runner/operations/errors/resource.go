package errors

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
	diffutils "github.com/kyverno/chainsaw/pkg/utils/diff"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type resourceError struct {
	expected unstructured.Unstructured
	actual   unstructured.Unstructured
	template bool
	bindings binding.Bindings
	errs     field.ErrorList
}

func ResourceError(
	expected unstructured.Unstructured,
	actual unstructured.Unstructured,
	template bool,
	bindings binding.Bindings,
	errs field.ErrorList,
) error {
	return resourceError{
		expected: expected,
		actual:   actual,
		template: template,
		bindings: bindings,
		errs:     errs,
	}
}

func (e resourceError) Error() string {
	var lines []string
	header := fmt.Sprintf("%s/%s/%s", e.actual.GetAPIVersion(), e.actual.GetKind(), client.Name(client.Key(&e.actual)))
	sep := strings.Repeat("-", len(header))
	lines = append(lines, sep, header, sep)
	if len(e.errs) != 0 {
		errLines := make([]string, 0, len(e.errs))
		for _, err := range e.errs {
			errLines = append(errLines, fmt.Sprintf("* %s", err))
		}
		slices.Sort(errLines)
		lines = append(lines, errLines...)
	}
	expected := e.expected
	var templateErr error
	if e.template {
		template := v1alpha1.Any{
			Value: expected.UnstructuredContent(),
		}
		if merged, err := templating.TemplateAndMerge(context.TODO(), expected, e.bindings, template); err != nil {
			templateErr = err
		} else {
			expected = merged
		}
	}
	if templateErr != nil {
		lines = append(lines, fmt.Sprintf("* ERROR: failed to compute expected template: %s", templateErr))
	}
	diff, err := diffutils.PrettyDiff(expected, *e.actual.DeepCopy())
	if err != nil {
		lines = append(lines, fmt.Sprintf("* %s", err))
	} else {
		lines = append(lines, "", diff)
	}
	return strings.Join(lines, "\n")
}
