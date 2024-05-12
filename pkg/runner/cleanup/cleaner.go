package cleanup

import (
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Cleaner = func(unstructured.Unstructured, clusters.Cluster)
