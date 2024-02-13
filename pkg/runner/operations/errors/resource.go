package errors

import (
	"fmt"
	"slices"
	"strings"

	"github.com/kyverno/chainsaw/pkg/client"
	diffutils "github.com/kyverno/chainsaw/pkg/utils/diff"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type resourceError struct {
	expected unstructured.Unstructured
	actual   unstructured.Unstructured
	errs     field.ErrorList
}

func ResourceError(
	expected unstructured.Unstructured,
	actual unstructured.Unstructured,
	errs field.ErrorList,
) error {
	return resourceError{
		expected: expected,
		actual:   actual,
		errs:     errs,
	}
}

func (e resourceError) Error() string {
	var lines []string
	header := fmt.Sprintf("%s/%s/%s", e.actual.GetAPIVersion(), e.actual.GetKind(), client.Name(client.ObjectKey(&e.actual)))
	sep := strings.Repeat("-", len(header))
	lines = append(lines, sep, header, sep)
	var errLines []string
	for _, err := range e.errs {
		errLines = append(errLines, fmt.Sprintf("* %s", err))
	}
	slices.Sort(errLines)
	lines = append(lines, errLines...)
	diff, err := diffutils.PrettyDiff(e.expected, *e.actual.DeepCopy())
	if err != nil {
		lines = append(lines, fmt.Sprintf("* %s", err))
	} else {
		lines = append(lines, "", diff)
	}
	return strings.Join(lines, "\n")
}
