package runner

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
// 	"github.com/kyverno/chainsaw/pkg/discovery"
// 	"github.com/stretchr/testify/assert"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/rest"
// 	tclock "k8s.io/utils/clock/testing"
// )

// type MockMainStart struct {
// 	code int
// }

// func (m *MockMainStart) Run() int {
// 	return m.code
// }

// func TestRun(t *testing.T) {
// 	fakeClock := tclock.NewFakePassiveClock(time.Now())

// 	tests := []struct {
// 		name       string
// 		tests      []discovery.Test
// 		config     v1alpha1.ConfigurationSpec
// 		restConfig *rest.Config
// 		mockReturn int
// 		wantErr    bool
// 	}{{
// 		name:       "Zero Tests",
// 		tests:      []discovery.Test{},
// 		config:     v1alpha1.ConfigurationSpec{},
// 		restConfig: &rest.Config{},
// 		wantErr:    false,
// 	}, {
// 		name: "Nil Rest Config with 1 Test",
// 		tests: []discovery.Test{
// 			{
// 				Err: nil,
// 				Test: &v1alpha1.Test{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name: "test1",
// 					},
// 				},
// 			},
// 		},
// 		config: v1alpha1.ConfigurationSpec{
// 			ReportFormat: v1alpha1.JSONFormat,
// 		},
// 		restConfig: nil,
// 		wantErr:    false,
// 	}, {
// 		name:  "Zero Tests with JSON Report",
// 		tests: []discovery.Test{},
// 		config: v1alpha1.ConfigurationSpec{
// 			ReportFormat: v1alpha1.JSONFormat,
// 		},
// 		restConfig: &rest.Config{},
// 		wantErr:    false,
// 	}, {
// 		name: "Success Case with 1 Test",
// 		tests: []discovery.Test{
// 			{
// 				Err: nil,
// 				Test: &v1alpha1.Test{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name: "test1",
// 					},
// 				},
// 			},
// 		},
// 		restConfig: &rest.Config{},
// 		mockReturn: 0,
// 		wantErr:    false,
// 	}, {
// 		name: "Failure Case with 1 Test",
// 		tests: []discovery.Test{
// 			{
// 				Err: nil,
// 				Test: &v1alpha1.Test{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name: "test1",
// 					},
// 				},
// 			},
// 		},
// 		restConfig: &rest.Config{},
// 		mockReturn: 2,
// 		wantErr:    true,
// 	}, {
// 		name: "Success Case with 1 Test with XML Report",
// 		tests: []discovery.Test{
// 			{
// 				Err: nil,
// 				Test: &v1alpha1.Test{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name: "test1",
// 					},
// 				},
// 			},
// 		},
// 		config: v1alpha1.ConfigurationSpec{
// 			ReportFormat: v1alpha1.XMLFormat,
// 			ReportName:   "chainsaw",
// 		},
// 		restConfig: &rest.Config{},
// 		mockReturn: 0,
// 		wantErr:    false,
// 	}, {
// 		name: "Error in saving Report",
// 		tests: []discovery.Test{
// 			{
// 				Err: nil,
// 				Test: &v1alpha1.Test{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name: "test1",
// 					},
// 				},
// 			},
// 		},
// 		config: v1alpha1.ConfigurationSpec{
// 			ReportFormat: "abc",
// 		},
// 		restConfig: &rest.Config{},
// 		mockReturn: 0,
// 		wantErr:    true,
// 	}}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockMainStart := &MockMainStart{
// 				code: tt.mockReturn,
// 			}
// 			_, err := run(context.TODO(), tt.restConfig, fakeClock, tt.config, mockMainStart, nil, tt.tests...)
// 			if tt.wantErr {
// 				assert.Error(t, err, "Run() should return an error")
// 			} else {
// 				assert.NoError(t, err, "Run() should not return an error")
// 			}
// 		})
// 	}
// }
