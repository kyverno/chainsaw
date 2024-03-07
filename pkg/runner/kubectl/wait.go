package kubectl

import (
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
)

func Wait(client client.Client, bindings binding.Bindings, collector *v1alpha1.Wait) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, errors.New("collector is null")
	}
	name, err := ConvertString(collector.Name, bindings)
	if err != nil {
		return nil, err
	}
	namespace, err := ConvertString(collector.Namespace, bindings)
	if err != nil {
		return nil, err
	}
	selector, err := ConvertString(collector.Selector, bindings)
	if err != nil {
		return nil, err
	}
	format, err := ConvertString(string(collector.Format), bindings)
	if err != nil {
		return nil, err
	}
	cluster, err := ConvertString(collector.Cluster, bindings)
	if err != nil {
		return nil, err
	}
	if name != "" && selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	resource, clustered, err := mapResource(client, bindings, collector.ResourceReference)
	if err != nil {
		return nil, err
	}
	cmd := v1alpha1.Command{
		Cluster:    cluster,
		Timeout:    collector.Timeout,
		Entrypoint: "kubectl",
		Args:       []string{"wait", resource},
	}
	if collector.For.Deletion != nil {
		cmd.Args = append(cmd.Args, "--for=delete")
	} else if collector.For.Condition != nil {
		name, err := ConvertString(collector.For.Condition.Name, bindings)
		if err != nil {
			return nil, err
		}
		if name == "" {
			return nil, errors.New("a condition name must be specified for condition wait type")
		}
		if collector.For.Condition.Value != nil {
			value, err := ConvertString(*collector.For.Condition.Value, bindings)
			if err != nil {
				return nil, err
			}
			cmd.Args = append(cmd.Args, fmt.Sprintf(`--for=condition=%s="%s"`, name, value))
		} else {
			cmd.Args = append(cmd.Args, fmt.Sprintf("--for=condition=%s", name))
		}
	} else {
		return nil, errors.New("either a deletion or a condition must be specified")
	}
	if name != "" {
		cmd.Args = append(cmd.Args, name)
	} else if selector != "" {
		cmd.Args = append(cmd.Args, "-l", selector)
	} else {
		cmd.Args = append(cmd.Args, "--all")
	}
	if !clustered {
		if namespace == "*" {
			cmd.Args = append(cmd.Args, "--all-namespaces")
		} else {
			if namespace == "" {
				namespace = "$NAMESPACE"
			}
			cmd.Args = append(cmd.Args, "-n", namespace)
		}
	}
	if format != "" {
		cmd.Args = append(cmd.Args, "-o", format)
	}
	// disable default timeout in the command
	cmd.Args = append(cmd.Args, "--timeout=-1s")
	return &cmd, nil
}
