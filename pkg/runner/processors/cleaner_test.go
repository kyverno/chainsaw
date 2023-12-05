package processors

import (
	"context"
	"testing"
	"time"

	fake "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	mock "github.com/kyverno/chainsaw/pkg/runner/operations/testing"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_Cleaner_Register(t *testing.T) {
	type registerTestCase struct {
		name       string
		timeout    time.Duration
		expectedOp int
	}

	testCases := []registerTestCase{
		{
			name:       "With 5 seconds timeout",
			timeout:    5 * time.Second,
			expectedOp: 1,
		},
		{
			name:       "With 10 seconds timeout",
			timeout:    10 * time.Second,
			expectedOp: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := &fake.FakeClient{}
			mockObj := unstructured.Unstructured{}
			fakeNamespacer := namespacer.New(fakeClient, "default")

			c := newCleaner(fakeNamespacer, nil)
			for i := 0; i < tc.expectedOp; i++ {
				localTimeout := tc.timeout
				c.register(mockObj, fakeClient, &localTimeout)
			}

			assert.Len(t, c.operations, tc.expectedOp)
			for _, op := range c.operations {
				assert.Equal(t, true, op.continueOnError)
				assert.GreaterOrEqual(t, *op.timeout, tc.timeout)
			}
		})
	}
}

func Test_Cleaner_Run(t *testing.T) {
	tests := []struct {
		name       string
		namespacer namespacer.Namespacer
		delay      *metav1.Duration
		operations []operation
	}{
		{
			name: "With 5 seconds delay",
			delay: &metav1.Duration{
				Duration: 5 * time.Second,
			},
			operations: []operation{
				{
					continueOnError: true,
					timeout:         nil,
					operation: mock.MockOperation{
						ExecFn: func(ctx context.Context) error {
							return nil
						},
					},
					operationReport: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cleaner{
				namespacer: tt.namespacer,
				delay:      tt.delay,
				operations: tt.operations,
			}
			c.run(context.Background())
		})
	}
}
