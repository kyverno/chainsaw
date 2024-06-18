package failer

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/testing"
)

var defaultFailer = New(false)

type Failer interface {
	Fail(context.Context, bool)
	FailNow(context.Context, bool)
}

type failer struct {
	pause bool
}

func New(pause bool) Failer {
	return failer{
		pause: pause,
	}
}

func (f failer) Fail(ctx context.Context, inconclusive bool) {
	f.wait()
	t := testing.FromContext(ctx)
	if !inconclusive {
		t.Fail()
	}
}

func (f failer) FailNow(ctx context.Context, inconclusive bool) {
	f.wait()
	t := testing.FromContext(ctx)
	if !inconclusive {
		t.FailNow()
	} else {
		t.SkipNow()
	}
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

func Fail(ctx context.Context, inconclusive bool) {
	f := getFailerOrDefault(ctx)
	f.Fail(ctx, inconclusive)
}

func FailNow(ctx context.Context, inconclusive bool) {
	f := getFailerOrDefault(ctx)
	f.FailNow(ctx, inconclusive)
}
