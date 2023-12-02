package validation

import (
	"fmt"
	"path/filepath"
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateOperation(t *testing.T) {
	exampleApply := &v1alpha1.Apply{
		FileRefOrResource: v1alpha1.FileRefOrResource{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
		},
	}
	exampleAssert := &v1alpha1.Assert{
		FileRefOrResource: v1alpha1.FileRefOrResource{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
		},
	}
	exampleCommand := &v1alpha1.Command{
		Entrypoint: "echo",
		Args:       []string{"hello world"},
	}
	exampleCreate := &v1alpha1.Create{
		FileRefOrResource: v1alpha1.FileRefOrResource{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
		},
	}
	exampleDelete := &v1alpha1.Delete{
		ObjectReference: v1alpha1.ObjectReference{
			APIVersion: "v1",
			Kind:       "Pod",
			ObjectSelector: v1alpha1.ObjectSelector{
				Namespace: "chainsaw",
				Labels: map[string]string{
					"app": "chainsaw",
				},
			},
		},
	}
	exampleError := &v1alpha1.Error{
		FileRefOrResource: v1alpha1.FileRefOrResource{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
		},
	}
	exampleScript := &v1alpha1.Script{
		Content: "echo 'hello world'",
	}
	tests := []struct {
		name      string
		input     v1alpha1.Operation
		expectErr bool
		errMsg    string
	}{{
		name:      "No operation statements provided",
		input:     v1alpha1.Operation{},
		expectErr: true,
		errMsg:    "no statement found in operation",
	}, {
		name: "Multiple operation statements provided",
		input: v1alpha1.Operation{
			Apply:   exampleApply,
			Assert:  exampleAssert,
			Command: exampleCommand,
		},
		expectErr: true,
		errMsg:    fmt.Sprintf("only one statement is allowed per operation (found %d)", 3),
	}, {
		name: "Only Apply operation statement provided",
		input: v1alpha1.Operation{
			Apply: exampleApply,
		},
		expectErr: false,
	}, {
		name: "Only Assert operation statement provided",
		input: v1alpha1.Operation{
			Assert: exampleAssert,
		},
		expectErr: false,
	}, {
		name: "Only Command operation statement provided",
		input: v1alpha1.Operation{
			Command: exampleCommand,
		},
		expectErr: false,
	}, {
		name: "Only Create operation statement provided",
		input: v1alpha1.Operation{
			Create: exampleCreate,
		},
		expectErr: false,
	}, {
		name: "Only Delete operation statement provided",
		input: v1alpha1.Operation{
			Delete: exampleDelete,
		},
		expectErr: false,
	}, {
		name: "Only Error operation statement provided",
		input: v1alpha1.Operation{
			Error: exampleError,
		},
		expectErr: false,
	}, {
		name: "Only Script operation statement provided",
		input: v1alpha1.Operation{
			Script: exampleScript,
		},
		expectErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateOperation(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
