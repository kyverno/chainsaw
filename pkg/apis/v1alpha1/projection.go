package v1alpha1

import (
	"encoding/json"

	"github.com/kyverno/kyverno-json/pkg/utils/copy"
)

// Projection can be any type.
// +k8s:deepcopy-gen=false
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:validation:Type:=""
type Projection struct {
	_value any
}

func NewProjection(value any) Projection {
	return Projection{
		_value: value,
	}
}

func (a *Projection) Value() any {
	return a._value
}

func (a *Projection) MarshalJSON() ([]byte, error) {
	return json.Marshal(a._value)
}

func (a *Projection) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	a._value = v
	return nil
}

func (in *Projection) DeepCopyInto(out *Projection) {
	out._value = copy.DeepCopy(in._value)
}

func (in *Projection) DeepCopy() *Projection {
	if in == nil {
		return nil
	}
	out := new(Projection)
	in.DeepCopyInto(out)
	return out
}
