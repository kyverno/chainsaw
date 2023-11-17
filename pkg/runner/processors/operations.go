package processors

// import (
// 	"path/filepath"

// 	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
// 	"github.com/kyverno/chainsaw/pkg/resource"
// 	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
// 	"github.com/kyverno/chainsaw/pkg/runner/logging"
// 	"github.com/kyverno/chainsaw/pkg/runner/operations"
// 	"github.com/kyverno/chainsaw/pkg/runner/operations/apply"
// 	"github.com/kyverno/chainsaw/pkg/runner/operations/assert"
// 	"github.com/kyverno/chainsaw/pkg/runner/operations/command"
// 	"github.com/kyverno/chainsaw/pkg/runner/operations/create"
// 	"github.com/kyverno/chainsaw/pkg/runner/operations/delete"
// 	"github.com/kyverno/chainsaw/pkg/runner/operations/error"
// 	"github.com/kyverno/chainsaw/pkg/runner/operations/script"
// 	"github.com/kyverno/kyverno/ext/output/color"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
// )

// func applyOperation(op v1alpha1.Apply) operations.Operation {
// 	var resources []ctrlclient.Object
// 	if operation.Create.Resource != nil {
// 		resources = append(resources, operation.Create.Resource)
// 	} else {
// 		loaded, err := resource.Load(filepath.Join(test.BasePath, operation.Create.File))
// 		if err != nil {
// 			logging.FromContext(ctx).Log("LOAD  ", color.BoldRed, err)
// 			fail(t, operation.ContinueOnError)
// 		}
// 		for i := range loaded {
// 			resources = append(resources, &loaded[i])
// 		}
// 	}

// 	return apply.New(c.getClient(dryRun), obj, getCleaner(cleanup, dryRun), check)
// }

// func assertOperation(expected unstructured.Unstructured) operations.Operation {
// 	return assert.New(c.client, expected)
// }

// func commandOperation(op v1alpha1.Command) operations.Operation {
// 	return command.New(op, c.namespace)
// }

// func createOperation(obj ctrlclient.Object, dryRun bool, check interface{}, cleanup cleanup.Cleaner) operations.Operation {
// 	return create.New(c.getClient(dryRun), obj, getCleaner(cleanup, dryRun), check)
// }

// func deleteOperation(operation *v1alpha1.Delete) operations.Operation {
// 	var resource unstructured.Unstructured
// 	resource.SetAPIVersion(operation.APIVersion)
// 	resource.SetKind(operation.Kind)
// 	resource.SetName(operation.Name)
// 	resource.SetNamespace(operation.Namespace)
// 	resource.SetLabels(operation.Labels)
// 	return delete.New(c.client, &resource)
// }

// func errorOperation(expected unstructured.Unstructured) operations.Operation {
// 	return error.New(c.client, expected)
// }

// func scriptOperation(exec v1alpha1.Script) operations.Operation {
// 	return script.New(exec, c.namespace)
// }
