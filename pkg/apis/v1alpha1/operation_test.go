package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperation_Bindings(t *testing.T) {
	type fields struct {
		Apply   *Apply
		Assert  *Assert
		Command *Command
		Create  *Create
		Delete  *Delete
		Error   *Error
		Patch   *Patch
		Script  *Script
		Sleep   *Sleep
		Update  *Update
		Wait    *Wait
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{{
		fields: fields{
			Apply: &Apply{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Assert: &Assert{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Command: &Command{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Create: &Create{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Delete: &Delete{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Error: &Error{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Patch: &Patch{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Script: &Script{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Sleep: &Sleep{},
		},
	}, {
		fields: fields{
			Update: &Update{
				ActionBindings: ActionBindings{Bindings: []Binding{{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Wait: &Wait{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Operation{
				Apply:   tt.fields.Apply,
				Assert:  tt.fields.Assert,
				Command: tt.fields.Command,
				Create:  tt.fields.Create,
				Delete:  tt.fields.Delete,
				Error:   tt.fields.Error,
				Patch:   tt.fields.Patch,
				Script:  tt.fields.Script,
				Sleep:   tt.fields.Sleep,
				Update:  tt.fields.Update,
				Wait:    tt.fields.Wait,
			}
			got := c.Bindings()
			assert.Equal(t, tt.want, len(got))
		})
	}
	assert.Panics(t, func() { (&Operation{}).Bindings() })
}

func TestOperation_Outputs(t *testing.T) {
	type fields struct {
		Apply   *Apply
		Assert  *Assert
		Command *Command
		Create  *Create
		Delete  *Delete
		Error   *Error
		Patch   *Patch
		Script  *Script
		Sleep   *Sleep
		Update  *Update
		Wait    *Wait
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{{
		fields: fields{
			Apply: &Apply{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", Any{Value: "bar"}}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Assert: &Assert{},
		},
	}, {
		fields: fields{
			Command: &Command{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", Any{Value: "bar"}}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Create: &Create{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", Any{Value: "bar"}}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Delete: &Delete{},
		},
	}, {
		fields: fields{
			Error: &Error{},
		},
	}, {
		fields: fields{
			Patch: &Patch{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", Any{Value: "bar"}}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Script: &Script{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", Any{Value: "bar"}}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Sleep: &Sleep{},
		},
	}, {
		fields: fields{
			Update: &Update{
				ActionOutputs: ActionOutputs{Outputs: []Output{{Binding: Binding{"foo", Any{Value: "bar"}}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Wait: &Wait{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Operation{
				Apply:   tt.fields.Apply,
				Assert:  tt.fields.Assert,
				Command: tt.fields.Command,
				Create:  tt.fields.Create,
				Delete:  tt.fields.Delete,
				Error:   tt.fields.Error,
				Patch:   tt.fields.Patch,
				Script:  tt.fields.Script,
				Sleep:   tt.fields.Sleep,
				Update:  tt.fields.Update,
				Wait:    tt.fields.Wait,
			}
			got := c.Outputs()
			assert.Equal(t, tt.want, len(got))
		})
	}
	assert.Panics(t, func() { (&Operation{}).Outputs() })
}
