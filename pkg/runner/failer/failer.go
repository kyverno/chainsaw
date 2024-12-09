package failer

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/testing"
)

var defaultFailer = New(false)

type Failer interface {
	Fail(context.Context, testing.TTest)
}

type failer struct {
	pause bool
}

func New(pause bool) Failer {
	return failer{
		pause: pause,
	}
}

func (f failer) Fail(ctx context.Context, t testing.TTest) {
	f.wait()
	t.Fail()
}

func (f failer) wait() {
	if f.pause {
		fmt.Println("Failure detected, press ENTER to continue...")
		fmt.Scanln() //nolint:errcheck
	}
}

func getFailerOrDefault(ctx context.Context) Failer {
	f := FromContext(ctx)
	if f == nil {
		return defaultFailer
	}
	return f
}

func Fail(ctx context.Context, t testing.TTest) {
	f := getFailerOrDefault(ctx)
	f.Fail(ctx, t)
}
