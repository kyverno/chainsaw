package model

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_summay(t *testing.T) {
	var wg sync.WaitGroup
	var s Summary
	const count int32 = 10000
	for i := 0; i < int(count); i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			s.IncFailed()
		}()
		go func() {
			defer wg.Done()
			s.IncPassed()
		}()
		go func() {
			defer wg.Done()
			s.IncSkipped()
		}()
	}
	wg.Wait()
	assert.Equal(t, count, s.Failed())
	assert.Equal(t, count, s.Passed())
	assert.Equal(t, count, s.Skipped())
}
