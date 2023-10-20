package operations

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/match"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

func Assert(ctx context.Context, expected unstructured.Unstructured, c client.Client) error {
	return wait.PollUntilContextCancel(ctx, interval, false, func(ctx context.Context) (bool, error) {
		candidates, err := read(ctx, expected, c)
		if err != nil {
			if errors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		if len(candidates) == 0 {
			return false, fmt.Errorf("no resource found")
		}

		var mismatchErrors []error
		matchedCount := 0
		for _, candidate := range candidates {
			err := match.Match(expected.UnstructuredContent(), candidate.UnstructuredContent())
			if err != nil {
				mismatchErrors = append(mismatchErrors, fmt.Errorf("resource %s/%s mismatch: %s", candidate.GetNamespace(), candidate.GetName(), err.Error()))
			} else {
				matchedCount++
			}
		}

		// If some matched and some not
		if matchedCount > 0 && matchedCount < len(candidates) {
			errMsg := fmt.Sprintf("%d of %d resources matched. Mismatches:", matchedCount, len(candidates))
			for _, err := range mismatchErrors {
				errMsg += "\n" + err.Error()
			}
			return false, fmt.Errorf(errMsg)
		}

		// If none of the resources matched.
		if matchedCount == 0 {
			errMsg := "All resources mismatched: "
			for _, err := range mismatchErrors {
				errMsg += "\n" + err.Error()
			}
			return false, fmt.Errorf(errMsg)
		}
		return true, nil
	})
}
