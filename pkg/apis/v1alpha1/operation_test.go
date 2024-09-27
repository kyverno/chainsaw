package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperation_Bindings(t *testing.T) {
	tests := []struct {
		name      string
		operation Operation
		want      int
	}{{
		operation: Operation{
			Apply: &Apply{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", NewProjection("bar")}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Assert: &Assert{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", NewProjection("bar")}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Command: &Command{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", NewProjection("bar")}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Create: &Create{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", NewProjection("bar")}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Delete: &Delete{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", NewProjection("bar")}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Describe: &Describe{},
		},
		want: 0,
	}, {
		operation: Operation{
			Error: &Error{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", NewProjection("bar")}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Events: &Events{},
		},
		want: 0,
	}, {
		operation: Operation{
			Get: &Get{},
		},
		want: 0,
	}, {
		operation: Operation{
			Patch: &Patch{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", NewProjection("bar")}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			PodLogs: &PodLogs{},
		},
		want: 0,
	}, {
		operation: Operation{
			Proxy: &Proxy{},
		},
		want: 0,
	}, {
		operation: Operation{
			Script: &Script{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", NewProjection("bar")}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Sleep: &Sleep{},
		},
	}, {
		operation: Operation{
			Update: &Update{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", NewProjection("bar")}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Wait: &Wait{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.operation.Bindings()
			assert.Equal(t, tt.want, len(got))
		})
	}
	assert.Panics(t, func() { (&Operation{}).Bindings() })
}

func TestOperation_Outputs(t *testing.T) {
	tests := []struct {
		name      string
		operation Operation
		want      int
	}{{
		operation: Operation{
			Apply: &Apply{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", NewProjection("bar")}}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Assert: &Assert{},
		},
	}, {
		operation: Operation{
			Command: &Command{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", NewProjection("bar")}}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Create: &Create{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", NewProjection("bar")}}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Delete: &Delete{},
		},
	}, {
		operation: Operation{
			Describe: &Describe{},
		},
	}, {
		operation: Operation{
			Error: &Error{},
		},
	}, {
		operation: Operation{
			Events: &Events{},
		},
	}, {
		operation: Operation{
			Get: &Get{},
		},
	}, {
		operation: Operation{
			Patch: &Patch{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", NewProjection("bar")}}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			PodLogs: &PodLogs{},
		},
	}, {
		operation: Operation{
			Proxy: &Proxy{},
		},
	}, {
		operation: Operation{
			Script: &Script{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", NewProjection("bar")}}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Sleep: &Sleep{},
		},
	}, {
		operation: Operation{
			Update: &Update{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", NewProjection("bar")}}}},
			},
		},
		want: 1,
	}, {
		operation: Operation{
			Wait: &Wait{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.operation.Outputs()
			assert.Equal(t, tt.want, len(got))
		})
	}
	assert.Panics(t, func() { (&Operation{}).Outputs() })
}
