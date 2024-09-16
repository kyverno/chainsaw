package kubectl

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
)

func Wait(ctx context.Context, client client.Client, tc binding.Bindings, collector *v1alpha1.Wait) (string, []string, error) {
	if collector == nil {
		return "", nil, errors.New("collector is null")
	}
	name, err := collector.Name.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	namespace, err := collector.Namespace.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	selector, err := collector.Selector.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	format, err := v1alpha1.Expression(collector.Format).Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	if name != "" && selector != "" {
		return "", nil, errors.New("name cannot be provided when a selector is specified")
	}
	resource, clustered, err := mapResource(ctx, client, tc, collector.ObjectType)
	if err != nil {
		return "", nil, err
	}
	args := []string{"wait", resource}
	if collector.WaitFor.Deletion != nil {
		args = append(args, "--for=delete")
	} else if collector.WaitFor.Condition != nil {
		name, err := collector.WaitFor.Condition.Name.Value(ctx, tc)
		if err != nil {
			return "", nil, err
		}
		if name == "" {
			return "", nil, errors.New("a condition name must be specified for condition wait type")
		}
		if collector.WaitFor.Condition.Value != nil {
			value, err := collector.WaitFor.Condition.Value.Value(ctx, tc)
			if err != nil {
				return "", nil, err
			}
			args = append(args, fmt.Sprintf(`--for=condition=%s=%s`, name, value))
		} else {
			args = append(args, fmt.Sprintf("--for=condition=%s", name))
		}
	} else if collector.WaitFor.JsonPath != nil {
		path, err := collector.WaitFor.JsonPath.Path.Value(ctx, tc)
		if err != nil {
			return "", nil, err
		}
		if path == "" {
			return "", nil, errors.New("a path must be specified for jsonpath wait type")
		}
		if collector.WaitFor.JsonPath.Value != nil {
			value, err := collector.WaitFor.JsonPath.Value.Value(ctx, tc)
			if err != nil {
				return "", nil, err
			}
			if value == "" {
				return "", nil, errors.New("a value must be specified for jsonpath wait type")
			}
			args = append(args, fmt.Sprintf(`--for=jsonpath=%s=%s`, path, value))
		} else {
			args = append(args, fmt.Sprintf(`--for=jsonpath=%s`, path))
		}
	} else {
		return "", nil, errors.New("either a deletion or a condition must be specified")
	}
	if name != "" {
		args = append(args, name)
	} else if selector != "" {
		args = append(args, "-l", selector)
	} else {
		args = append(args, "--all")
	}
	if !clustered {
		if namespace == "*" {
			args = append(args, "--all-namespaces")
		} else {
			if namespace == "" {
				namespace = "$NAMESPACE"
			}
			args = append(args, "-n", namespace)
		}
	}
	if format != "" {
		args = append(args, "-o", format)
	}
	if collector.Timeout != nil {
		args = append(args, "--timeout", collector.Timeout.Duration.String())
	} else {
		args = append(args, "--timeout=-1s")
	}
	return "kubectl", args, nil
}
