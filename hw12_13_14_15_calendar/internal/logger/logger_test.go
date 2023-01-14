package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name        string
		level       string
		funcName    string
		message     string
		expectedMsg string
	}{
		{
			name:        "error msg",
			level:       "ERROR",
			funcName:    "Error",
			message:     "This is error message",
			expectedMsg: "ERROR:This is error message",
		},
		{
			name:        "skipp_warn",
			level:       "ERROR",
			funcName:    "Warn",
			message:     "This is error message",
			expectedMsg: "",
		},
		{
			name:        "skipp_info",
			level:       "ERROR",
			funcName:    "Info",
			message:     "This is error message",
			expectedMsg: "",
		},
		{
			name:        "skipp_debug",
			level:       "ERROR",
			funcName:    "Debug",
			message:     "This is error message",
			expectedMsg: "",
		},
		{
			name:        "debug",
			level:       "DEBUG",
			funcName:    "Debug",
			message:     "This is error message",
			expectedMsg: "DEBUG:This is error message",
		},
		{
			name:        "skipp_debug2",
			level:       "INFO",
			funcName:    "Debug",
			message:     "This is error message",
			expectedMsg: "",
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("case %s", tc.name), func(t *testing.T) {
			var b bytes.Buffer
			tc := tc
			t.Parallel()
			l := New(tc.level, &b, nil)

			switch tc.funcName {
			case "Error":
				l.Errorf(tc.message)
			case "Warn":
				l.Warningf(tc.message)
			case "Info":
				l.Infof(tc.message)
			case "Debug":
				l.Debugf(tc.message)
			}

			require.Equal(t, tc.expectedMsg, b.String(), "error output message")
		})
	}
}
