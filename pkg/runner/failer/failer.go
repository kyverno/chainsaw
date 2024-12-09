package failer

import (
	"fmt"
)

var Default = New(false)

type Failer interface {
	Fail()
}

type failer struct {
	pause bool
}

func New(pause bool) Failer {
	return failer{
		pause: pause,
	}
}

func (f failer) Fail() {
	f.wait()
}

func (f failer) wait() {
	if f.pause {
		fmt.Println("Failure detected, press ENTER to continue...")
		fmt.Scanln() //nolint:errcheck
	}
}
