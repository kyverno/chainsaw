package kubectl

import (
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func WaitForResource(waiter *v1alpha1.Wait) (*v1alpha1.Command, error) {
	if waiter == nil {
		return nil, nil
	}
	if waiter.Resource == "" {
		return nil, errors.New("a resource must be specified")
	}
	if waiter.Condition == "" {
		return nil, errors.New("a condition must be specified")
	}
	if waiter.Name != "" && waiter.Selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}

	args := []string{"wait"}

	if waiter.Name != "" {
		args = append(args, fmt.Sprintf("%s/%s", waiter.Resource, waiter.Name))
	} else {
		args = append(args, waiter.Resource)
	}
	args = append(args, fmt.Sprintf("--for=condition=%s", waiter.Condition))

	if waiter.Selector != "" {
		args = append(args, "-l", waiter.Selector)
	}
	if waiter.AllNamespaces {
		args = append(args, "--all-namespaces")
	} else if waiter.Namespace != "" {
		args = append(args, "-n", waiter.Namespace)
	} else {
		args = append(args, "-n", "$NAMESPACE")
	}
	if waiter.Timeout != nil {
		args = append(args, fmt.Sprintf("--timeout=%s", waiter.Timeout.Duration.String()))
	}
	cmd := v1alpha1.Command{
		Cluster:    waiter.Cluster,
		Timeout:    waiter.Timeout,
		Entrypoint: "kubectl",
		Args:       args,
	}
	return &cmd, nil
}
