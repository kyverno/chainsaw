package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestDeps_MatchString(t *testing.T) {
	tests := []struct {
		name    string
		Test    bool
		pat     string
		str     string
		want    bool
		wantErr bool
	}{{
		name:    "empty pattern",
		Test:    false,
		pat:     "",
		str:     "foo",
		want:    true,
		wantErr: false,
	}, {
		name:    "empty string",
		Test:    false,
		pat:     "foo",
		str:     "",
		want:    false,
		wantErr: false,
	}, {
		name:    "both empty",
		Test:    false,
		pat:     "",
		str:     "",
		want:    true,
		wantErr: false,
	}, {
		name:    "match",
		Test:    false,
		pat:     "o w",
		str:     "hello world",
		want:    true,
		wantErr: false,
	}, {
		name:    "no match",
		Test:    false,
		pat:     "cat",
		str:     "hello world",
		want:    false,
		wantErr: false,
	}, {
		name:    "error",
		Test:    false,
		pat:     "^[$",
		str:     "hello world",
		want:    false,
		wantErr: true,
	}, {
		name:    "test",
		Test:    true,
		pat:     "^[$",
		str:     "hello world",
		want:    true,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &TestDeps{
				Test: tt.Test,
			}
			got, err := d.MatchString(tt.pat, tt.str)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestTestDeps_StartCPUProfile(t *testing.T) {
	d := &TestDeps{}
	assert.NoError(t, d.StartCPUProfile(nil))
}

func TestTestDeps_StopCPUProfile(t *testing.T) {
	d := &TestDeps{}
	d.StopCPUProfile()
}

func TestTestDeps_WriteProfileTo(t *testing.T) {
	d := &TestDeps{}
	assert.NoError(t, d.WriteProfileTo("", nil, 0))
}

func TestTestDeps_ImportPath(t *testing.T) {
	d := &TestDeps{}
	assert.Equal(t, "", d.ImportPath())
}

func TestTestDeps_StartTestLog(t *testing.T) {
	d := &TestDeps{}
	d.StartTestLog(nil)
}

func TestTestDeps_StopTestLog(t *testing.T) {
	d := &TestDeps{}
	assert.NoError(t, d.StopTestLog())
}

func TestTestDeps_CoordinateFuzzing(t *testing.T) {
	d := &TestDeps{}
	err := d.CoordinateFuzzing(0, 0, 0, 0, 0, nil, nil, "", "")
	assert.NoError(t, err)
}

func TestTestDeps_RunFuzzWorker(t *testing.T) {
	d := &TestDeps{}
	err := d.RunFuzzWorker(nil)
	assert.NoError(t, err)
}

func TestTestDeps_ReadCorpus(t *testing.T) {
	d := &TestDeps{}
	got, err := d.ReadCorpus("", nil)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

func TestTestDeps_CheckCorpus(t *testing.T) {
	d := &TestDeps{}
	assert.NoError(t, d.CheckCorpus(nil, nil))
}

func TestTestDeps_ResetCoverage(t *testing.T) {
	d := &TestDeps{}
	d.ResetCoverage()
}

func TestTestDeps_SnapshotCoverage(t *testing.T) {
	d := &TestDeps{}
	d.SnapshotCoverage()
}

func TestTestDeps_InitRuntimeCoverage(t *testing.T) {
	d := &TestDeps{}
	mode, tearDown, snapcov := d.InitRuntimeCoverage()
	assert.Empty(t, mode)
	assert.Nil(t, tearDown)
	assert.Nil(t, snapcov)
}
