package test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateOperation(t *testing.T) {
	exampleApply := &v1alpha1.Apply{
		ActionResourceRef: v1alpha1.ActionResourceRef{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
		},
	}
	exampleAssert := &v1alpha1.Assert{
		ActionCheckRef: v1alpha1.ActionCheckRef{
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
		ActionResourceRef: v1alpha1.ActionResourceRef{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
		},
	}
	exampleDelete := &v1alpha1.Delete{
		Ref: &v1alpha1.ObjectReference{
			ObjectType: v1alpha1.ObjectType{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectName: v1alpha1.ObjectName{
				Namespace: "chainsaw",
			},
			Labels: map[string]string{
				"app": "chainsaw",
			},
		},
	}
	exampleError := &v1alpha1.Error{
		ActionCheckRef: v1alpha1.ActionCheckRef{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
		},
	}
	examplePatch := &v1alpha1.Patch{
		ActionResourceRef: v1alpha1.ActionResourceRef{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
		},
	}
	exampleScript := &v1alpha1.Script{
		Content: "echo 'hello world'",
	}
	exampleSleep := &v1alpha1.Sleep{
		Duration: metav1.Duration{Duration: 5 * time.Second},
	}
	exampleUpdate := &v1alpha1.Update{
		ActionResourceRef: v1alpha1.ActionResourceRef{
			FileRef: v1alpha1.FileRef{
				File: filepath.Join("..", "..", "testdata", "validation", "example-file.yaml"),
			},
		},
	}
	exampleWait := &v1alpha1.Wait{
		ActionObject: v1alpha1.ActionObject{
			ObjectType: v1alpha1.ObjectType{
				APIVersion: "v1",
				Kind:       "Pod",
			},
		},
		WaitFor: v1alpha1.WaitFor{
			Deletion: &v1alpha1.WaitForDeletion{},
		},
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
		name: "Only Patch operation statement provided",
		input: v1alpha1.Operation{
			Patch: examplePatch,
		},
		expectErr: false,
	}, {
		name: "Only Script operation statement provided",
		input: v1alpha1.Operation{
			Script: exampleScript,
		},
		expectErr: false,
	}, {
		name: "Only Sleep operation statement provided",
		input: v1alpha1.Operation{
			Sleep: exampleSleep,
		},
		expectErr: false,
	}, {
		name: "Only Update operation statement provided",
		input: v1alpha1.Operation{
			Update: exampleUpdate,
		},
		expectErr: false,
	}, {
		name: "Only Wait operation statement provided",
		input: v1alpha1.Operation{
			Wait: exampleWait,
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
