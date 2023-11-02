package testing

import (
	"fmt"
)

type FakeLogger struct {
	Messages []string
}

func (tl *FakeLogger) Log(args ...interface{}) {
	for _, arg := range args {
		tl.Messages = append(tl.Messages, fmt.Sprint(arg))
	}
}

func (tl *FakeLogger) Helper() {}
