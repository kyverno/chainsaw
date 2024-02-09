package validation

import (
	"fmt"
	"testing"
	"time"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

func TestValidateCatch(t *testing.T) {
	exampleEvents := &v1alpha1.Events{}
	examplePodLogs := &v1alpha1.PodLogs{
		Selector: "app=example",
		Tail:     ptr.To(10),
	}
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
		Resource: "pods",
	}
	tests := []struct {
		name      string
		input     v1alpha1.Catch
		expectErr bool
		errMsg    string
	}{{
		name:      "No Catch statements provided",
		input:     v1alpha1.Catch{},
		expectErr: true,
		errMsg:    "no statement found in operation",
	}, {
		name: "Multiple Catch statements provided",
		input: v1alpha1.Catch{
			PodLogs: examplePodLogs,
			Events:  exampleEvents,
			Command: exampleCommand,
		},
		expectErr: true,
		errMsg:    fmt.Sprintf("only one statement is allowed per operation (found %d)", 3),
	}, {
		name: "Only PodLogs statement provided",
		input: v1alpha1.Catch{
			PodLogs: examplePodLogs,
		},
		expectErr: false,
	}, {
		name: "Only Events statement provided",
		input: v1alpha1.Catch{
			Events: exampleEvents,
		},
		expectErr: false,
	}, {
		name: "Only Command statement provided",
		input: v1alpha1.Catch{
			Command: exampleCommand,
		},
		expectErr: false,
	}, {
		name: "Only Script statement provided",
		input: v1alpha1.Catch{
			Script: exampleScript,
		},
		expectErr: false,
	}, {
		name: "Only Sleep statement provided",
		input: v1alpha1.Catch{
			Sleep: exampleSleep,
		},
		expectErr: false,
	}, {
		name: "Only Describe statement provided",
		input: v1alpha1.Catch{
			Describe: exampleDescribe,
		},
		expectErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateCatch(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
