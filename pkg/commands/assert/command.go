package assert

import (
	"context"
	"fmt"
	"io"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	ctrlClient "github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/resource"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	nspacer "github.com/kyverno/chainsaw/pkg/runner/namespacer"
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
	filePath            string
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
	cmd.Flags().StringVarP(&opts.filePath, "file", "f", "", "Path to the file to assert or '-' to read from stdin")
	cmd.Flags().DurationVar(&opts.timeout.Duration, "timeout", v1alpha1.DefaultAssertTimeout, "The assert timeout to use")
	cmd.Flags().StringVar(&opts.namespace, "namespace", "default", "Namespace to use")
	cmd.Flags().BoolVar(&opts.noColor, "no-color", false, "Removes output colors")
	clientcmd.BindOverrideFlags(&opts.kubeConfigOverrides, cmd.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return cmd
}

func preRunE(opts *options, cmd *cobra.Command, args []string) error {
	if len(args) > 0 && opts.filePath == "" {
		opts.filePath = args[0]
	} else if opts.filePath == "" {
		return fmt.Errorf("either a file path as an argument or the --file flag must be provided")
	}
	return nil
}

func runE(opts options, cmd *cobra.Command, client ctrlClient.Client, namespacer nspacer.Namespacer) error {
	color.Init(opts.noColor, true)
	out := cmd.OutOrStdout()
	var resources []unstructured.Unstructured
	var err error

	if opts.filePath == "-" {
		content, readErr := io.ReadAll(cmd.InOrStdin())
		if readErr != nil {
			return fmt.Errorf("failed to read from stdin: %w", readErr)
		}
		resources, err = resource.Parse(content, false)
		if err != nil {
			return fmt.Errorf("failed to parse stdin content: %w", err)
		}
	} else {
		resources, err = resource.Load(opts.filePath, false)
		if err != nil {
			return fmt.Errorf("failed to load file '%s': %w", opts.filePath, err)
		}
	}

	if client == nil {
		cfg, err := restutils.Config(opts.kubeConfigOverrides)
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
	for _, resource := range resources {
		if err := assert(opts, client, resource, namespacer); err != nil {
			return fmt.Errorf("assertion failed: %w", err)
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
