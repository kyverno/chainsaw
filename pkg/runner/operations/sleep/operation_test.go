package sleep

import (
	"context"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	ttesting "github.com/kyverno/chainsaw/pkg/testing"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_operation_Exec(t *testing.T) {
	tests := []struct {
		name         string
		sleep        v1alpha1.Sleep
		expectedLogs []string
	}{{
		name:         "zero",
		sleep:        v1alpha1.Sleep{},
		expectedLogs: []string{"SLEEP: RUN - []", "SLEEP: DONE - []"},
	}, {
		name: "1s",
		sleep: v1alpha1.Sleep{
			Duration: metav1.Duration{Duration: time.Second},
		},
		expectedLogs: []string{"SLEEP: RUN - []", "SLEEP: DONE - []"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			operation := New(
				tt.sleep,
			)
			logger := &tlogging.FakeLogger{}
			err := operation.Exec(ttesting.IntoContext(logging.IntoContext(ctx, logger), t), nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedLogs, logger.Logs)
		})
	}
}
