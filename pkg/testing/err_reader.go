package testing

import (
	"errors"
)

type ErrReader struct{}

func (e *ErrReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading from stdin")
}
