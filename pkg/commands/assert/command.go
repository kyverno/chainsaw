package assert

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/resource"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	opassert "github.com/kyverno/chainsaw/pkg/runner/operations/assert"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/kyverno/kyverno/ext/output/color"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/clientcmd"
)

type options struct {
	timeout             metav1.Duration
	namespace           string
	noColor             bool
	kubeConfigOverrides clientcmd.ConfigOverrides
}

func Command() *cobra.Command {
	var options options
	cmd := &cobra.Command{
		Use:          "assert",
		Short:        "Evaluate assertion",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			color.Init(options.noColor, true)
			out := cmd.OutOrStdout()
			file := args[0]
			resources, err := resource.Load(file, false)
			if err != nil {
				return err
			}
			cfg, err := restutils.Config(options.kubeConfigOverrides)
			if err != nil {
				return err
			}
			client, err := client.New(cfg)
			if err != nil {
				return err
			}
			client = runnerclient.New(client)
			namespacer := namespacer.New(client, options.namespace)
			for _, resource := range resources {
				if err := assert(options, client, resource, namespacer); err != nil {
					return fmt.Errorf("Assertion FAILED: %w", err)
				}
			}
			fmt.Fprintln(out, "Assertion(s) PASSED")
			return nil
		},
	}
	cmd.Flags().DurationVar(&options.timeout.Duration, "timeout", v1alpha1.DefaultAssertTimeout, "The assert timeout to use")
	cmd.Flags().StringVar(&options.namespace, "namespace", "default", "Namespace to use")
	cmd.Flags().BoolVar(&options.noColor, "no-color", false, "Removes output colors")
	clientcmd.BindOverrideFlags(&options.kubeConfigOverrides, cmd.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return cmd
}

func assert(options options, client client.Client, resource unstructured.Unstructured, namespacer namespacer.Namespacer) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.timeout.Duration)
	defer cancel()
	op := opassert.New(client, resource, namespacer, binding.NewBindings(), false)
	return op.Exec(ctx)
}
