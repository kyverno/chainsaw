package test

import (
	"fmt"
	"testing"
	"time"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateCatchFinally(t *testing.T) {
	examplePodLogs := &v1alpha1.PodLogs{
		ActionObjectSelector: v1alpha1.ActionObjectSelector{
			Selector: "app=example",
		},
	}
	exampleEvents := &v1alpha1.Events{}
	exampleCommand := &v1alpha1.Command{
		Entrypoint: "echo",
		Args:       []string{"Hello, World!"},
	}
	exampleScript := &v1alpha1.Script{
		Content: "echo Hello, World!",
	}
	exampleSleep := &v1alpha1.Sleep{
		Duration: metav1.Duration{Duration: 5 * time.Second},
	}
	exampleDescribe := &v1alpha1.Describe{
		ActionObject: v1alpha1.ActionObject{
			ObjectType: v1alpha1.ObjectType{
				APIVersion: "v1",
				Kind:       "Pod",
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
		For: v1alpha1.For{
			Condition: &v1alpha1.Condition{
				Name: "Ready",
			},
		},
	}
	exampleGet := &v1alpha1.Get{
		ActionObject: v1alpha1.ActionObject{
			ObjectType: v1alpha1.ObjectType{
				APIVersion: "v1",
				Kind:       "Pod",
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
	tests := []struct {
		name      string
		input     v1alpha1.CatchFinally
		expectErr bool
		errMsg    string
	}{{
		name:      "No Finally statements provided",
		input:     v1alpha1.CatchFinally{},
		expectErr: true,
		errMsg:    "no statement found in operation",
	}, {
		name: "Multiple Finally statements provided",
		input: v1alpha1.CatchFinally{
			PodLogs: examplePodLogs,
			Events:  exampleEvents,
			Command: exampleCommand,
		},
		expectErr: true,
		errMsg:    fmt.Sprintf("only one statement is allowed per operation (found %d)", 3),
	}, {
		name: "Only PodLogs statement provided",
		input: v1alpha1.CatchFinally{
			PodLogs: examplePodLogs,
		},
		expectErr: false,
	}, {
		name: "Only Events statement provided",
		input: v1alpha1.CatchFinally{
			Events: exampleEvents,
		},
		expectErr: false,
	}, {
		name: "Only Command statement provided",
		input: v1alpha1.CatchFinally{
			Command: exampleCommand,
		},
		expectErr: false,
	}, {
		name: "Only Script statement provided",
		input: v1alpha1.CatchFinally{
			Script: exampleScript,
		},
		expectErr: false,
	}, {
		name: "Only Sleep statement provided",
		input: v1alpha1.CatchFinally{
			Sleep: exampleSleep,
		},
		expectErr: false,
	}, {
		name: "Only Describe statement provided",
		input: v1alpha1.CatchFinally{
			Describe: exampleDescribe,
		},
		expectErr: false,
	}, {
		name: "Only Wait statement provided",
		input: v1alpha1.CatchFinally{
			Wait: exampleWait,
		},
		expectErr: false,
	}, {
		name: "Only Get statement provided",
		input: v1alpha1.CatchFinally{
			Get: exampleGet,
		},
		expectErr: false,
	}, {
		name: "Only Delete statement provided",
		input: v1alpha1.CatchFinally{
			Delete: exampleDelete,
		},
		expectErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateCatchFinally(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
