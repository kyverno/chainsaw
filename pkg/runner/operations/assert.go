package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/match"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

func Assert(ctx context.Context, expected unstructured.Unstructured, c client.Client) error {
	return wait.PollUntilContextCancel(ctx, interval, false, func(ctx context.Context) (bool, error) {
		if candidates, err := read(ctx, expected, c); err != nil {
			if errors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		} else {
			for _, candidate := range candidates {
				// at least one must match
				if err := match.Match(expected.UnstructuredContent(), candidate.UnstructuredContent()); err == nil {
					return true, nil
				}
			}
		}
		return false, nil
	})
}
