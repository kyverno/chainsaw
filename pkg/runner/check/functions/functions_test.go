package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestGetFunctions(t *testing.T) {
	var c ctrlclient.Client
	assert.Equal(t, 1, len(GetFunctions(c)))
}
