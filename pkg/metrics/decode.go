package metrics

import (
	"errors"
	"io"
	"sort"
	"strings"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)

func Decode(in string, ts model.Time) (model.Vector, error) {
	dec := &expfmt.SampleDecoder{
		Dec: expfmt.NewDecoder(strings.NewReader(in), expfmt.NewFormat(expfmt.TypeTextPlain)),
		Opts: &expfmt.DecodeOptions{
			Timestamp: ts,
		},
	}
	var all model.Vector
	for {
		var smpls model.Vector
		err := dec.Decode(&smpls)
		if err != nil && errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		all = append(all, smpls...)
	}
	sort.Sort(all)
	return all, nil
}
