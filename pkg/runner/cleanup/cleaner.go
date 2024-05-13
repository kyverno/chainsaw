package cleanup

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Cleaner = func(unstructured.Unstructured, client.Client)
