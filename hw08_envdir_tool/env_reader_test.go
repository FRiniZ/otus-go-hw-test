package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		envr Environment
		err  error
	}{
		{
			name: "simple test",
			dir:  "./testdata/env/",
			envr: Environment{
				"BAR":   {"bar", false},
				"EMPTY": {"", true},
				"FOO":   {"   foo\nwith new line", false},
				"HELLO": {"\"hello\"", false},
				"UNSET": {"", true},
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			envr, err := ReadDir(tc.dir)
			require.Truef(t, reflect.DeepEqual(envr, tc.envr), "", err)
		})
	}
}
