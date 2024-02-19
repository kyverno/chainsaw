package kubectl

import (
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func WaitForResource(waiter *v1alpha1.Wait) (*v1alpha1.Command, error) {
	if waiter == nil {
		return nil, errors.New("waiter configuration cannot be nil")
	}
	if waiter.Resource == "" {
		return nil, errors.New("a resource must be specified")
	}
	if waiter.Condition == "" {
		return nil, errors.New("a condition must be specified")
	}

	cmd := v1alpha1.Command{
		Cluster:    waiter.Cluster,
		Entrypoint: "kubectl",
		Args:       []string{"wait", waiter.Resource},
	}

	if waiter.ResourceName != "" {
		cmd.Args = append(cmd.Args, waiter.ResourceName)
	}

	if waiter.ObjectLabelsSelector.Selector != "" {
		cmd.Args = append(cmd.Args, "-l", waiter.ObjectLabelsSelector.Selector)
	}

	if waiter.Timeout != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--timeout=%s", waiter.Timeout.String()))
	}

	if waiter.IncludeUninitialized != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--include-uninitialized=%t", *waiter.IncludeUninitialized))
	}

	cmd.Args = append(cmd.Args, fmt.Sprintf("--for=condition=%s", waiter.Condition))

	if waiter.PollInterval != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--poll-interval=%s", waiter.PollInterval.String()))
	}

	return &cmd, nil
}
