package operations

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/match"
	"github.com/pmezard/go-difflib/difflib"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

func Assert(ctx context.Context, expected unstructured.Unstructured, c client.Client) error {
	var differences []string
	var successfulMatches int
	var totalCandidatesChecked int
	var candidates []unstructured.Unstructured
	printedDifferences := make(map[string]bool) // This will help us track which resources have had their differences printed

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
					resourceString := fmt.Sprintf("%v", candidate.UnstructuredContent())
					if _, exists := printedDifferences[resourceString]; !exists {
						differences = append(differences, getDifference(expected.UnstructuredContent(), candidate.UnstructuredContent()))
						printedDifferences[resourceString] = true
					}
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
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Context timeout. Successful matches: %d. Remaining unchecked candidates: %d.\n",
			successfulMatches, len(candidates)-totalCandidatesChecked))

		if successfulMatches == 0 {
			for _, diff := range differences {
				sb.WriteString(diff)
			}
		}
		return fmt.Errorf(sb.String())
	}
	return err
}

func getDifference(expectedContent, candidateContent map[string]interface{}) string {
	expectedContentBytes, _ := yaml.Marshal(expectedContent)
	candidateContentBytes, _ := yaml.Marshal(candidateContent)

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(expectedContentBytes)),
		B:        difflib.SplitLines(string(candidateContentBytes)),
		FromFile: "Expected",
		ToFile:   "Actual",
		Context:  3,
	}
	diffStr, _ := difflib.GetUnifiedDiffString(diff)
	return diffStr
}
