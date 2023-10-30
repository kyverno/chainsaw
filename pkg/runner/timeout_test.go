package runner

import (
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func Test_timeout(t *testing.T) {
	tests := []struct {
		name   string
		config v1alpha1.ConfigurationSpec
		test   v1alpha1.TestSpec
		step   v1alpha1.TestStepSpec
		want   *time.Duration
	}{{
		name: "none",
		config: v1alpha1.ConfigurationSpec{
			Timeout: nil,
		},
		test: v1alpha1.TestSpec{
			Timeout: nil,
		},
		step: v1alpha1.TestStepSpec{
			Timeout: nil,
		},
		want: nil,
	}, {
		name: "from config",
		config: v1alpha1.ConfigurationSpec{
			Timeout: &metav1.Duration{Duration: 1 * time.Minute},
		},
		test: v1alpha1.TestSpec{
			Timeout: nil,
		},
		step: v1alpha1.TestStepSpec{
			Timeout: nil,
		},
		want: ptr.To(1 * time.Minute),
	}, {
		name: "from test",
		config: v1alpha1.ConfigurationSpec{
			Timeout: &metav1.Duration{Duration: 1 * time.Minute},
		},
		test: v1alpha1.TestSpec{
			Timeout: &metav1.Duration{Duration: 2 * time.Minute},
		},
		step: v1alpha1.TestStepSpec{
			Timeout: nil,
		},
		want: ptr.To(2 * time.Minute),
	}, {
		name: "from step",
		config: v1alpha1.ConfigurationSpec{
			Timeout: &metav1.Duration{Duration: 1 * time.Minute},
		},
		test: v1alpha1.TestSpec{
			Timeout: &metav1.Duration{Duration: 2 * time.Minute},
		},
		step: v1alpha1.TestStepSpec{
			Timeout: &metav1.Duration{Duration: 3 * time.Minute},
		},
		want: ptr.To(3 * time.Minute),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timeout(tt.config, tt.test, tt.step)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cancelNoOp(t *testing.T) {
	tests := []struct {
		name string
	}{{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cancelNoOp()
		})
	}
}

func Test_timeoutCtx(t *testing.T) {
	tests := []struct {
		name   string
		config v1alpha1.ConfigurationSpec
		test   v1alpha1.TestSpec
		step   v1alpha1.TestStepSpec
		want   bool
	}{{
		name: "none",
		config: v1alpha1.ConfigurationSpec{
			Timeout: nil,
		},
		test: v1alpha1.TestSpec{
			Timeout: nil,
		},
		step: v1alpha1.TestStepSpec{
			Timeout: nil,
		},
		want: false,
	}, {
		name: "from config",
		config: v1alpha1.ConfigurationSpec{
			Timeout: &metav1.Duration{Duration: 1 * time.Minute},
		},
		test: v1alpha1.TestSpec{
			Timeout: nil,
		},
		step: v1alpha1.TestStepSpec{
			Timeout: nil,
		},
		want: true,
	}, {
		name: "from test",
		config: v1alpha1.ConfigurationSpec{
			Timeout: nil,
		},
		test: v1alpha1.TestSpec{
			Timeout: &metav1.Duration{Duration: 2 * time.Minute},
		},
		step: v1alpha1.TestStepSpec{
			Timeout: nil,
		},
		want: true,
	}, {
		name: "from step",
		config: v1alpha1.ConfigurationSpec{
			Timeout: nil,
		},
		test: v1alpha1.TestSpec{
			Timeout: nil,
		},
		step: v1alpha1.TestStepSpec{
			Timeout: &metav1.Duration{Duration: 3 * time.Minute},
		},
		want: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, cancel := timeoutCtx(tt.config, tt.test, tt.step)
			defer cancel()
			assert.NotNil(t, got)
			assert.NotNil(t, cancel)
			_, ok := got.Deadline()
			assert.Equal(t, tt.want, ok)
		})
	}
}

func Test_timeoutExecCtx(t *testing.T) {
	tests := []struct {
		name   string
		config v1alpha1.ConfigurationSpec
		test   v1alpha1.TestSpec
		step   v1alpha1.TestStepSpec
		exec   v1alpha1.Exec
		want   bool
	}{{
		name: "none",
		config: v1alpha1.ConfigurationSpec{
			Timeout: nil,
		},
		test: v1alpha1.TestSpec{
			Timeout: nil,
		},
		step: v1alpha1.TestStepSpec{
			Timeout: nil,
		},
		exec: v1alpha1.Exec{
			Timeout: nil,
		},
		want: false,
	}, {
		name: "from config",
		config: v1alpha1.ConfigurationSpec{
			Timeout: &metav1.Duration{Duration: 1 * time.Minute},
		},
		test: v1alpha1.TestSpec{
			Timeout: nil,
		},
		step: v1alpha1.TestStepSpec{
			Timeout: nil,
		},
		exec: v1alpha1.Exec{
			Timeout: nil,
		},
		want: true,
	}, {
		name: "from test",
		config: v1alpha1.ConfigurationSpec{
			Timeout: nil,
		},
		test: v1alpha1.TestSpec{
			Timeout: &metav1.Duration{Duration: 2 * time.Minute},
		},
		step: v1alpha1.TestStepSpec{
			Timeout: nil,
		},
		exec: v1alpha1.Exec{
			Timeout: nil,
		},
		want: true,
	}, {
		name: "from step",
		config: v1alpha1.ConfigurationSpec{
			Timeout: nil,
		},
		test: v1alpha1.TestSpec{
			Timeout: nil,
		},
		step: v1alpha1.TestStepSpec{
			Timeout: &metav1.Duration{Duration: 3 * time.Minute},
		},
		exec: v1alpha1.Exec{
			Timeout: nil,
		},
		want: true,
	}, {
		name: "from exec",
		config: v1alpha1.ConfigurationSpec{
			Timeout: nil,
		},
		test: v1alpha1.TestSpec{
			Timeout: nil,
		},
		step: v1alpha1.TestStepSpec{
			Timeout: nil,
		},
		exec: v1alpha1.Exec{
			Timeout: &metav1.Duration{Duration: 4 * time.Minute},
		},
		want: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, cancel := timeoutExecCtx(tt.exec, tt.config, tt.test, tt.step)
			defer cancel()
			assert.NotNil(t, got)
			assert.NotNil(t, cancel)
			_, ok := got.Deadline()
			assert.Equal(t, tt.want, ok)
		})
	}
}
