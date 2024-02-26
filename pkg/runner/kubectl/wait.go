package kubectl

import (
	"errors"
	"fmt"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func Wait(collector *v1alpha1.Wait) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, nil
	}
	if collector.Resource == "" {
		return nil, errors.New("a resource must be specified")
	}
	if collector.Name != "" && collector.Selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}

	args := []string{"wait"}

	if collector.Name != "" {
		args = append(args, fmt.Sprintf("%s/%s", collector.Resource, collector.Name))
	} else {
		args = append(args, collector.Resource)
	}

	if collector.For.Deletion != nil {
		args = append(args, "--for=delete")
	} else if collector.For.Condition != nil {
		if collector.For.Condition.ConditionName == "" {
			return nil, errors.New("a condition name must be specified for condition wait type")
		}
		if collector.For.Condition.ConditionValue != "" {
			args = append(args, fmt.Sprintf("--for=condition=%s=%s", collector.For.Condition.ConditionName, collector.For.Condition.ConditionValue))
		} else {
			args = append(args, fmt.Sprintf("--for=condition=%s", collector.For.Condition.ConditionName))
		}
	} else {
		return nil, errors.New("either a deletion or a condition must be specified")
	}

	if collector.Selector != "" {
		args = append(args, "-l", collector.Selector)
	}

	if collector.Namespace != "" {
		args = append(args, "-n", collector.Namespace)
	} else {
		args = append(args, "-n", "$NAMESPACE")
	}
	if collector.Timeout != nil {
		args = append(args, fmt.Sprintf("--timeout=%s", collector.Timeout.Duration.String()))
	} else {
		args = append(args, fmt.Sprintf("--timeout=%s", (time.Duration(-1)*time.Second).String()))
	}
	if collector.OutputFormat != nil {
		args = append(args, fmt.Sprintf("--output=%s", *collector.OutputFormat))
	}

	cmd := v1alpha1.Command{
		Cluster:    collector.Cluster,
		Timeout:    collector.Timeout,
		Entrypoint: "kubectl",
		Args:       args,
	}
	return &cmd, nil
}
