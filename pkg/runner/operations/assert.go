package operations

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/match"
	"github.com/pmezard/go-difflib/difflib"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

func Assert(ctx context.Context, expected unstructured.Unstructured, c client.Client) error {
	var differences []string
	var successfulMatches int
	var totalCandidatesChecked int
	var candidates []unstructured.Unstructured
	err := wait.PollUntilContextCancel(ctx, interval, false, func(ctx context.Context) (bool, error) {
		var err error
		if candidates, err = read(ctx, expected, c); err != nil {
			if errors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		} else {
			for _, candidate := range candidates {
				totalCandidatesChecked++
				if err := match.Match(expected.UnstructuredContent(), candidate.UnstructuredContent()); err != nil {
					// Store the difference for this failed match
					differences = append(differences, getDifference(expected.UnstructuredContent(), candidate.UnstructuredContent()))
				} else {
					// at least one match found
					successfulMatches++
					return true, nil
				}
			}
		}
		return false, nil
	})

	// Handle context timeout
	if err != nil && ctx.Err() == context.DeadlineExceeded {
		fmt.Printf("Context timeout. Successful matches: %d. Remaining unchecked candidates: %d.\n",
			successfulMatches, len(candidates)-totalCandidatesChecked)

		if successfulMatches == 0 {
			for _, diff := range differences {
				fmt.Println(diff)
			}
		}
	}
	return err
}

func getDifference(expectedContent, candidateContent map[string]interface{}) string {
	expectedContentStr := fmt.Sprintf("%v", expectedContent)
	candidateContentStr := fmt.Sprintf("%v", candidateContent)

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expectedContentStr),
		B:        difflib.SplitLines(candidateContentStr),
		FromFile: "Expected",
		ToFile:   "Candidate",
		Context:  3,
	}
	diffStr, _ := difflib.GetUnifiedDiffString(diff)
	return diffStr
}
