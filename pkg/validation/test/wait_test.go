package test

import (
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateWait(t *testing.T) {
	tests := []struct {
		name      string
		input     *v1alpha1.Wait
		expectErr bool
		errMsg    string
	}{{
		name:      "No resource provided",
		input:     &v1alpha1.Wait{},
		expectErr: true,
		errMsg:    "kind must be specified",
	}, {
		name:      "No resource provided",
		input:     &v1alpha1.Wait{},
		expectErr: true,
		errMsg:    "apiVersion must be specified",
	}, {
		name: "No condition provided",
		input: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
			},
		},
		expectErr: true,
		errMsg:    "either a deletion, condition or a jsonpath must be specified",
	}, {
		name: "Neither Name nor Selector provided",
		input: &v1alpha1.Wait{
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
		},
		expectErr: false,
	}, {
		name: "Both Name and Selector provided",
		input: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "example-name",
					},
					Selector: "foo=bar",
				},
			},
			For: v1alpha1.For{
				Condition: &v1alpha1.Condition{
					Name: "Ready",
				},
			},
		},
		expectErr: true,
		errMsg:    "a name or label selector must be specified (found both)",
	}, {
		name: "Only Name provided",
		input: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					ObjectName: v1alpha1.ObjectName{
						Name: "example-name",
					},
				},
			},
			For: v1alpha1.For{
				Deletion: &v1alpha1.Deletion{},
			},
		},
		expectErr: false,
	}, {
		name: "Only Selector provided",
		input: &v1alpha1.Wait{
			ActionObject: v1alpha1.ActionObject{
				ObjectType: v1alpha1.ObjectType{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ActionObjectSelector: v1alpha1.ActionObjectSelector{
					Selector: "example-selector",
				},
			},
			For: v1alpha1.For{
				Deletion: &v1alpha1.Deletion{},
			},
		},
		expectErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateWait(field.NewPath("wait"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
