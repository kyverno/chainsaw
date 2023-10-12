package runner

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/resource"
	"k8s.io/client-go/rest"
)

func Run(cfg *rest.Config, tests ...v1alpha1.Test) (int, error) {
	if len(tests) == 0 {
		return 0, nil
	}
	flag.Parse()
	testing.Init()
	// Set the verbose test flag to true since we are not using the regular go test CLI.
	if err := flag.Set("test.v", "true"); err != nil {
		return 0, err
	}
	// TODO: flags
	var testDeps testDeps
	m := testing.MainStart(
		&testDeps,
		[]testing.InternalTest{{
			Name: "chainsaw",
			F: func(t *testing.T) {
				t.Helper()
				run(t, cfg, tests...)
			},
		}},
		nil,
		nil,
		nil,
	)
	return m.Run(), nil
}

func run(t *testing.T, cfg *rest.Config, tests ...v1alpha1.Test) {
	t.Helper()
	for _, test := range tests {
		client, err := client.New(cfg)
		if err != nil {
			t.Fatal(err)
		}
		func(t *testing.T, test v1alpha1.Test) {
			t.Helper()
			t.Run(test.Name, func(t *testing.T) {
				t.Helper()
				t.Parallel()
				for i, step := range test.Spec.Steps {
					func(t *testing.T, step v1alpha1.TestStepSpec) {
						t.Helper()
						t.Run(fmt.Sprintf("step-%d", i+1), func(t *testing.T) {
							t.Helper()
							executeStep(t, step, client)
						})
					}(t, step)
				}
			})
		}(t, test)
	}
}

func executeStep(t *testing.T, step v1alpha1.TestStepSpec, client client.Client) {
	t.Helper()
	for _, apply := range step.Apply {
		resources, err := resource.Load(apply.File)
		if err != nil {
			t.Fatal(err)
		}
		for i := range resources {
			_, err := client.Apply(context.Background(), &resources[i])
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
