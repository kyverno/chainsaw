package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatch_Bindings(t *testing.T) {
	type fields struct {
		PodLogs  *PodLogs
		Events   *Events
		Describe *Describe
		Wait     *Wait
		Get      *Get
		Delete   *Delete
		Command  *Command
		Script   *Script
		Sleep    *Sleep
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{{
		fields: fields{
			Command: &Command{
				Bindings: []Binding{{"foo", Any{Value: "bar"}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Delete: &Delete{
				Bindings: []Binding{{"foo", Any{Value: "bar"}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Describe: &Describe{},
		},
	}, {
		fields: fields{
			Events: &Events{},
		},
	}, {
		fields: fields{
			Get: &Get{},
		},
	}, {
		fields: fields{
			PodLogs: &PodLogs{},
		},
	}, {
		fields: fields{
			Script: &Script{
				Bindings: []Binding{{"foo", Any{Value: "bar"}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Sleep: &Sleep{},
		},
	}, {
		fields: fields{
			Wait: &Wait{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catch{
				PodLogs:  tt.fields.PodLogs,
				Events:   tt.fields.Events,
				Describe: tt.fields.Describe,
				Wait:     tt.fields.Wait,
				Get:      tt.fields.Get,
				Delete:   tt.fields.Delete,
				Command:  tt.fields.Command,
				Script:   tt.fields.Script,
				Sleep:    tt.fields.Sleep,
			}
			got := c.Bindings()
			assert.Equal(t, tt.want, len(got))
		})
	}
	assert.Panics(t, func() { (&Catch{}).Bindings() })
}

func TestCatch_Outputs(t *testing.T) {
	type fields struct {
		PodLogs  *PodLogs
		Events   *Events
		Describe *Describe
		Wait     *Wait
		Get      *Get
		Delete   *Delete
		Command  *Command
		Script   *Script
		Sleep    *Sleep
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{{
		fields: fields{
			Command: &Command{
				Outputs: []Output{{Binding: Binding{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Delete: &Delete{},
		},
	}, {
		fields: fields{
			Describe: &Describe{},
		},
	}, {
		fields: fields{
			Events: &Events{},
		},
	}, {
		fields: fields{
			Get: &Get{},
		},
	}, {
		fields: fields{
			PodLogs: &PodLogs{},
		},
	}, {
		fields: fields{
			Script: &Script{
				Outputs: []Output{{Binding: Binding{"foo", Any{Value: "bar"}}}},
			},
		},
		want: 1,
	}, {
		fields: fields{
			Sleep: &Sleep{},
		},
	}, {
		fields: fields{
			Wait: &Wait{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catch{
				PodLogs:  tt.fields.PodLogs,
				Events:   tt.fields.Events,
				Describe: tt.fields.Describe,
				Wait:     tt.fields.Wait,
				Get:      tt.fields.Get,
				Delete:   tt.fields.Delete,
				Command:  tt.fields.Command,
				Script:   tt.fields.Script,
				Sleep:    tt.fields.Sleep,
			}
			got := c.Outputs()
			assert.Equal(t, tt.want, len(got))
		})
	}
	assert.Panics(t, func() { (&Catch{}).Outputs() })
}
