package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name string
		cmds []string
		envr Environment
		rc   int
	}{
		{
			name: "simple test",
			cmds: []string{"ls", "-ll"},
			rc:   0,
			envr: Environment{},
		},
		{
			name: "environment test",
			cmds: []string{"./testdata/env.sh"},
			rc:   0,
			envr: Environment{"FOO": {"foo", false}},
		},
		{
			name: "negative test",
			cmds: []string{"llls", "-ll"},
			rc:   -1,
			envr: Environment{},
		},
		{
			name: "negative test2",
			cmds: []string{"./testdata/err.sh"},
			rc:   10,
			envr: Environment{},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rc := RunCmd(tc.cmds, tc.envr)
			require.Truef(t, rc == tc.rc, "", nil)
		})
	}
}
