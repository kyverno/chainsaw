package functions

import (
	"github.com/kyverno/chainsaw/pkg/metrics"
	"github.com/prometheus/common/model"
)

func jpMetricsDecode(arguments []any) (any, error) {
	var text string
	if err := getArg(arguments, 0, &text); err != nil {
		return nil, err
	}
	vector, err := metrics.Decode(text, model.Now())
	if err != nil {
		return nil, err
	}
	return vector, nil
}
