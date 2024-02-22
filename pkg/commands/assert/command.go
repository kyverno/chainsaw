package assert

import (
	"context"
	"fmt"
	"io"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	ctrlClient "github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/resource"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	nspacer "github.com/kyverno/chainsaw/pkg/runner/namespacer"
	opassert "github.com/kyverno/chainsaw/pkg/runner/operations/assert"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/kyverno/kyverno/ext/output/color"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type options struct {
	timeout             metav1.Duration
	namespace           string
	noColor             bool
	assertPath          string
	resourcePath        string
	kubeConfigOverrides clientcmd.ConfigOverrides
}

func Command() *cobra.Command {
	var opts options
	cmd := &cobra.Command{
		Use:          "assert [flags] [FILE]",
		Short:        "Evaluate assertion",
		Args:         cobra.RangeArgs(0, 1),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunE(&opts, cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(opts, cmd, nil, nil)
		},
	}
	cmd.Flags().StringVarP(&opts.assertPath, "file", "f", "", "Path to the file to assert or '-' to read from stdin")
	cmd.Flags().StringVarP(&opts.resourcePath, "resource", "r", "", "Path to the file containing the resource")
	cmd.Flags().DurationVar(&opts.timeout.Duration, "timeout", v1alpha1.DefaultAssertTimeout, "The assert timeout to use")
	cmd.Flags().StringVar(&opts.namespace, "namespace", "default", "Namespace to use")
	cmd.Flags().BoolVar(&opts.noColor, "no-color", false, "Removes output colors")
	cmd.Flags().BoolVar(&opts.noColor, "clustered", false, "Defines if the resource is clustered (only applies when resource is loaded from a file)")
	clientcmd.BindOverrideFlags(&opts.kubeConfigOverrides, cmd.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return cmd
}

func preRunE(opts *options, cmd *cobra.Command, args []string) error {
	if len(args) > 0 && opts.assertPath == "" {
		opts.assertPath = args[0]
	} else if opts.assertPath == "" {
		return fmt.Errorf("either a file path as an argument or the --file flag must be provided")
	}
	return nil
}

func runE(opts options, cmd *cobra.Command, client ctrlClient.Client, namespacer nspacer.Namespacer) error {
	color.Init(opts.noColor, true)
	out := cmd.OutOrStdout()
	var assertions []unstructured.Unstructured
	var err error
	if opts.assertPath == "-" {
		content, readErr := io.ReadAll(cmd.InOrStdin())
		if readErr != nil {
			return fmt.Errorf("failed to read from stdin: %w", readErr)
		}
		assertions, err = resource.Parse(content, false)
		if err != nil {
			return fmt.Errorf("failed to parse stdin content: %w", err)
		}
	} else {
		assertions, err = resource.Load(opts.assertPath, false)
		if err != nil {
			return fmt.Errorf("failed to load file '%s': %w", opts.assertPath, err)
		}
	}
	if opts.resourcePath != "" {
		ressources, err := resource.Load(opts.resourcePath, false)
		if err != nil {
			return fmt.Errorf("failed to load file '%s': %w", opts.resourcePath, err)
		}
		client = &tclient.FakeClient{
			GetFn: func(_ context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				// TODO: we should improve the lookup logic here
				*obj.(*unstructured.Unstructured) = ressources[0]
				return nil
			},
			ListFn: func(_ context.Context, _ int, list ctrlclient.ObjectList, _ ...ctrlclient.ListOption) error {
				*list.(*unstructured.UnstructuredList) = unstructured.UnstructuredList{
					Items: ressources,
				}
				return nil
			},
			IsObjectNamespacedFn: func(_ int, _ runtime.Object) (bool, error) {
				return false, nil
			},
		}
	}
	if client == nil {
		cfg, err := restutils.DefaultConfig(opts.kubeConfigOverrides)
		if err != nil {
			return err
		}
		newClient, err := ctrlClient.New(cfg)
		if err != nil {
			return err
		}
		client = runnerclient.New(newClient)
	}
	if namespacer == nil {
		namespacer = nspacer.New(client, opts.namespace)
	}
	for _, assertion := range assertions {
		if err := assert(opts, client, assertion, namespacer); err != nil {
			return fmt.Errorf("assertion failed\n%w", err)
		}
	}
	fmt.Fprintln(out, "Assertion(s) PASSED")
	return nil
}

func assert(opts options, client ctrlClient.Client, resource unstructured.Unstructured, namespacer nspacer.Namespacer) error {
	ctx, cancel := context.WithTimeout(context.Background(), opts.timeout.Duration)
	defer cancel()
	op := opassert.New(client, resource, namespacer, binding.NewBindings(), false)
	return op.Exec(ctx)
}
