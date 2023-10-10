package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *v1alpha1.Configuration
		wantErr bool
	}{{
		name:    "confimap",
		path:    "../../testdata/config/configmap.yaml",
		wantErr: true,
	}, {
		name:    "not found",
		path:    "../../testdata/config/not-found.yaml",
		wantErr: true,
	}, {
		name:    "empty",
		path:    "../../testdata/config/empty.yaml",
		wantErr: true,
	}, {
		name: "default",
		path: "../../testdata/config/default.yaml",
		want: &v1alpha1.Configuration{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "chainsaw.kyverno.io/v1alpha1",
				Kind:       "Configuration",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
			Spec: v1alpha1.ConfigurationSpec{
				Timeout: &metav1.Duration{
					Duration: 30 * time.Second,
				},
			},
		},
	}, {
		name: "timeout 1m",
		path: "../../testdata/config/timeout-1m.yaml",
		want: &v1alpha1.Configuration{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "chainsaw.kyverno.io/v1alpha1",
				Kind:       "Configuration",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "timeout-1m",
			},
			Spec: v1alpha1.ConfigurationSpec{
				Timeout: &metav1.Duration{
					Duration: time.Minute,
				},
			},
		},
	}, {
		name:    "multiple",
		path:    "../../testdata/config/multiple.yaml",
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
