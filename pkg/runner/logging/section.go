package logging

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"go.uber.org/multierr"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

type section struct {
	name string
	args []any
}

func (s section) String() string {
	return strings.TrimSpace(s.name + "\n" + fmt.Sprint(s.args...))
}

func Section(name string, args ...any) fmt.Stringer {
	return section{
		name: "=== " + strings.ToUpper(name),
		args: args,
	}
}

func ErrSection(err error) fmt.Stringer {
	var errs []string
	for _, err := range multierr.Errors(err) {
		var agg utilerrors.Aggregate
		if errors.As(err, &agg) {
			for _, err := range agg.Errors() {
				errs = append(errs, err.Error())
			}
		} else {
			errs = append(errs, err.Error())
		}
	}
	slices.Sort(errs)
	return Section("ERROR", strings.Join(errs, "\n"))
}
