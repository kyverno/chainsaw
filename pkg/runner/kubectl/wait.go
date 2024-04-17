package kubectl

import (
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
)

func Wait(client client.Client, bindings binding.Bindings, collector *v1alpha1.Wait) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, errors.New("collector is null")
	}
	name, err := apibindings.String(collector.Name, bindings)
	if err != nil {
		return nil, err
	}
	namespace, err := apibindings.String(collector.Namespace, bindings)
	if err != nil {
		return nil, err
	}
	selector, err := apibindings.String(collector.Selector, bindings)
	if err != nil {
		return nil, err
	}
	format, err := apibindings.String(string(collector.Format), bindings)
	if err != nil {
		return nil, err
	}
	cluster, err := apibindings.String(collector.Cluster, bindings)
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
		name, err := apibindings.String(collector.For.Condition.Name, bindings)
		if err != nil {
			return nil, err
		}
		if name == "" {
			return nil, errors.New("a condition name must be specified for condition wait type")
		}
		if collector.For.Condition.Value != nil {
			value, err := apibindings.String(*collector.For.Condition.Value, bindings)
			if err != nil {
				return nil, err
			}
			cmd.Args = append(cmd.Args, fmt.Sprintf(`--for=condition=%s=%s`, name, value))
		} else {
			cmd.Args = append(cmd.Args, fmt.Sprintf("--for=condition=%s", name))
		}
	} else if collector.For.JsonPath != nil {
		path, err := apibindings.String(collector.For.JsonPath.Path, bindings)
		if err != nil {
			return nil, err
		}
		if path == "" {
			return nil, errors.New("a path must be specified for jsonpath wait type")
		}
		value, err := apibindings.String(collector.For.JsonPath.Value, bindings)
		if err != nil {
			return nil, err
		}
		if value == "" {
			return nil, errors.New("a value must be specified for jsonpath wait type")
		}
		cmd.Args = append(cmd.Args, fmt.Sprintf(`--for=jsonpath='%s'=%s`, path, value))
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
	if collector.Timeout != nil {
		cmd.Args = append(cmd.Args, "--timeout", collector.Timeout.Duration.String())
	} else {
		cmd.Args = append(cmd.Args, "--timeout=-1s")
	}
	return &cmd, nil
}
