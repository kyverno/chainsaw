package kubebuilder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func scaffoldCommand() *cobra.Command {
	var save bool
	var kind string
	var apiVersion string
	var description bool

	cmd := &cobra.Command{
		Use:   "scaffold [test-name]",
		Short: "Scaffold a Chainsaw test for a kubebuilder operator",
		Long: `Scaffold generates a chainsaw-test.yaml tailored for kubebuilder-based operators.
It creates a test that applies a custom resource, asserts its reconciled state,
and cleans up afterward.`,
		Example: `  # Print scaffolded test to stdout
  chainsaw kubebuilder scaffold mytest --kind MyKind --api-version mygroup.io/v1alpha1

  # Save directly into the test directory
  chainsaw kubebuilder scaffold mytest --kind MyKind --api-version mygroup.io/v1alpha1 --save`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			testName := strings.ToLower(strings.ReplaceAll(args[0], "_", "-"))

			test := buildTest(testName, kind, apiVersion, description)

			data, err := yaml.Marshal(&test)
			if err != nil {
				return fmt.Errorf("failed to marshal test: %w", err)
			}

			if save {
				dir := filepath.Join(".", args[0])
				if err := os.MkdirAll(dir, 0o755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", dir, err)
				}
				outPath := filepath.Join(dir, "chainsaw-test.yaml")
				fmt.Fprintf(out, "Saving file %s ...\n", outPath)
				if err := os.WriteFile(outPath, data, 0o600); err != nil {
					return fmt.Errorf("failed to write file: %w", err)
				}
			} else {
				fmt.Fprintln(out, string(data))
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&save, "save", false, "If set, saves the scaffolded test to a file instead of printing")
	cmd.Flags().StringVar(&kind, "kind", "MyKind", "The Kind of the custom resource under test")
	cmd.Flags().StringVar(&apiVersion, "api-version", "mygroup.io/v1alpha1", "The APIVersion of the custom resource under test")
	cmd.Flags().BoolVar(&description, "description", true, "If set, adds descriptions to the scaffolded test")

	return cmd
}

func buildTest(name, kind, apiVersion string, description bool) v1alpha1.Test {
	test := v1alpha1.Test{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.GroupVersion.String(),
			Kind:       "Test",
		},
		Spec: v1alpha1.TestSpec{
			Description: desc(description, fmt.Sprintf("e2e test for kubebuilder operator: %s", name)),
			Steps:       buildSteps(kind, apiVersion, description),
		},
	}
	test.SetName(name)
	return test
}

func buildSteps(kind, apiVersion string, description bool) []v1alpha1.TestStep {
	return []v1alpha1.TestStep{
		{
			Name: "create-resource",
			TestStepSpec: v1alpha1.TestStepSpec{
				Description: desc(description, fmt.Sprintf("Apply the %s custom resource and assert it is reconciled", kind)),
				Try: []v1alpha1.Operation{
					{
						OperationBase: v1alpha1.OperationBase{
							Description: desc(description, fmt.Sprintf("Apply the %s resource", kind)),
						},
						Apply: &v1alpha1.Apply{
							ActionResourceRef: v1alpha1.ActionResourceRef{
								FileRef: v1alpha1.FileRef{
									File: "resources.yaml",
								},
							},
						},
					},
					{
						OperationBase: v1alpha1.OperationBase{
							Description: desc(description, fmt.Sprintf("Assert that %s is reconciled successfully", kind)),
						},
						Assert: &v1alpha1.Assert{
							ActionCheckRef: v1alpha1.ActionCheckRef{
								FileRef: v1alpha1.FileRef{
									File: "assert.yaml",
								},
							},
						},
					},
				},
				Catch: []v1alpha1.CatchFinally{
					{
						Description: desc(description, "Collect events on failure"),
						Events: &v1alpha1.Events{
							ActionObjectSelector: v1alpha1.ActionObjectSelector{
								ObjectName: v1alpha1.ObjectName{
									Name: "sample-resource",
								},
							},
						},
					},
					{
						Description: desc(description, "Collect pod logs on failure"),
						PodLogs: &v1alpha1.PodLogs{
							ActionObjectSelector: v1alpha1.ActionObjectSelector{
								Selector: "control-plane=controller-manager",
							},
						},
					},
				},
				Finally: []v1alpha1.CatchFinally{
					{
						Description: desc(description, "Wait before cleanup"),
						Sleep: &v1alpha1.Sleep{
							Duration: metav1.Duration{Duration: 5 * time.Second},
						},
					},
				},
			},
		},
		{
			Name: "delete-resource",
			TestStepSpec: v1alpha1.TestStepSpec{
				Description: desc(description, fmt.Sprintf("Delete the %s custom resource and verify cleanup", kind)),
				Try: []v1alpha1.Operation{
					{
						OperationBase: v1alpha1.OperationBase{
							Description: desc(description, fmt.Sprintf("Delete the %s resource", kind)),
						},
						Delete: &v1alpha1.Delete{
							Ref: &v1alpha1.ObjectReference{
								ObjectType: v1alpha1.ObjectType{
									APIVersion: v1alpha1.Expression(apiVersion),
									Kind:       v1alpha1.Expression(kind),
								},
								ObjectName: v1alpha1.ObjectName{
									Name: "sample-resource",
								},
							},
						},
					},
				},
			},
		},
	}
}

func desc(enabled bool, s string) string {
	if !enabled {
		return ""
	}
	return s
}
